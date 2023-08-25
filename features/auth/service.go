package auth

import (
	"context"
	"errors"
	"fmt"
	"gobit-demo/model"
	"gobit-demo/util"

	"gorm.io/gorm"
)

var (
	ErrUserNotFound         = errors.New("用户不存在")
	ErrUserDisabled         = errors.New("用户已禁用")
	ErrDuplicatedUser       = errors.New("用户已存在")
	ErrIncorrectCredentials = errors.New("用户名或密码错误")
)

type Service interface {
	CreateUser(ctx context.Context, data *CreateUserRequest) (*model.User, error)
	FindUserByUserName(ctx context.Context, name string) (*model.User, error)
	FindUserByMobile(ctx context.Context, mobile string) (*model.User, error)
	Login(ctx context.Context, data *LoginRequest) (*LoginUser, error)
	GetUserPermissions(ctx context.Context, uid model.UUID) ([]string, error)
	GetUserRoles(ctx context.Context, uid model.UUID) ([]string, error)
}

type service struct {
	g *gorm.DB
	h PasswordHash
}

// GetUserPermissions implements Service.
func (s *service) GetUserPermissions(ctx context.Context, uid model.UUID) ([]string, error) {
	var permissions []string
	if err := s.g.WithContext(ctx).Model(&model.Permission{}).
		Distinct("permissions.name").
		Joins("INNER JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Joins("INNER JOIN user_roles ON user_roles.role_id = role_permissions.role_id").
		Where("user_roles.user_id = ?", uid).
		Pluck("permissions.name", &permissions).Error; err != nil {
		return nil, fmt.Errorf("get user permissions: %w", err)
	}
	return permissions, nil
}

// GetUserRoles implements Service.
func (s *service) GetUserRoles(ctx context.Context, uid model.UUID) ([]string, error) {
	var roles []string
	if err := s.g.WithContext(ctx).Model(&model.Role{}).
		Joins("INNER JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", uid).
		Pluck("roles.name", &roles).Error; err != nil {
		return nil, fmt.Errorf("get user roles: %w", err)
	}
	return roles, nil
}

func (s *service) FindUserByUserName(ctx context.Context, username string) (*model.User, error) {
	user := new(model.User)
	err := s.g.WithContext(ctx).Model(&model.User{}).Where(&model.User{Username: username}).Take(user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find user by name: %w", err)
	}
	return user, err
}

func (s *service) FindUserByMobile(ctx context.Context, mobile string) (*model.User, error) {
	user := new(model.User)
	err := s.g.WithContext(ctx).Model(&model.User{}).Where(&model.User{Mobile: mobile}).Take(user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find user by mobile: %w", err)
	}
	return user, err
}

func (s *service) CreateUser(ctx context.Context, data *CreateUserRequest) (*model.User, error) {
	user := &model.User{
		Name:     data.Name,
		Username: data.Username,
		Mobile:   data.Mobile,
	}
	if err := s.g.Transaction(func(tx *gorm.DB) error {
		exists, err := util.GormUtil.Exists(
			tx.WithContext(ctx).Model(&model.User{}).
				Where(&model.User{Mobile: data.Mobile}).
				Or(&model.User{Username: data.Username}),
		)
		if err != nil {
			return fmt.Errorf("check duplicate user: %w", err)
		}
		if exists {
			return ErrDuplicatedUser
		}

		h, err := s.h.Hash(data.Password)
		if err != nil {
			return fmt.Errorf("hash password: %w", err)
		}
		user.Password = h

		if err = tx.WithContext(ctx).Create(user).Error; err != nil {
			return fmt.Errorf("create user: %w", err)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *service) Login(ctx context.Context, data *LoginRequest) (*LoginUser, error) {
	u := new(LoginUser)
	if err := s.g.Transaction(func(tx *gorm.DB) error {
		m := new(model.User)
		err := s.g.WithContext(ctx).Where(&model.User{Mobile: data.Mobile}).
			Or(&model.User{Username: data.Username}).Take(m).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrIncorrectCredentials
		}
		if err != nil {
			return fmt.Errorf("find user: %w", err)
		}

		if m.Disabled {
			return ErrUserDisabled
		}

		if err := s.h.VerifyPassword(data.Password, m.Password); err != nil {
			return ErrIncorrectCredentials
		}

		p, err := s.GetUserPermissions(ctx, m.ID)
		if err != nil {
			return fmt.Errorf("get user permissions: %w", err)
		}

		r, err := s.GetUserRoles(ctx, m.ID)
		if err != nil {
			return fmt.Errorf("get user roles: %w", err)
		}

		u.Permissions = p
		u.Roles = r
		u.ID = m.ID
		u.Username = m.Username

		return nil
	}); err != nil {
		return nil, err
	}
	return u, nil
}

func NewService(g *gorm.DB) Service {
	return &service{g: g, h: NewBcryptPasswordHash()}
}

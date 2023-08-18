package auth

import (
	"context"
	"errors"
	"fmt"
	"gobit-demo/internal/orm"
	"gobit-demo/model"

	"gorm.io/gorm"
)

var (
	ErrUserNotFound         = errors.New("用户不存在")
	ErrUserDisabled         = errors.New("用户已禁用")
	ErrDuplicatedUser       = errors.New("用户已存在")
	ErrIncorrectCredentials = errors.New("用户名或密码错误")
)

type Service interface {
	CreateUser(ctx context.Context, data *CreateUserRequest) (*RegisterUser, error)
	FindUserByUserName(ctx context.Context, name string) (*LoginUser, error)
	FindUserByMobile(ctx context.Context, mobile string) (*LoginUser, error)
	FindUserByCredential(ctx context.Context, data *LoginRequest) (*LoginUser, error)
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

func (s *service) FindUserByUserName(ctx context.Context, username string) (*LoginUser, error) {
	user := new(LoginUser)
	err := s.g.WithContext(ctx).Model(&model.User{}).Where(&model.User{Username: username}).Take(user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find user by name: %w", err)
	}
	return user, err
}

func (s *service) FindUserByMobile(ctx context.Context, mobile string) (*LoginUser, error) {
	user := new(LoginUser)
	err := s.g.WithContext(ctx).Model(&model.User{}).Where(&model.User{Mobile: mobile}).Take(user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find user by mobile: %w", err)
	}
	return user, err
}

func (s *service) CreateUser(ctx context.Context, data *CreateUserRequest) (*RegisterUser, error) {
	user := model.User{
		Name:     data.Name,
		Username: data.Username,
		Mobile:   data.Mobile,
	}
	if err := s.g.Transaction(func(tx *gorm.DB) error {
		exists, err := orm.NewGormWrapper(tx.WithContext(ctx).Model(&model.User{}).
			Where(&model.User{Mobile: data.Mobile}).
			Or(&model.User{Username: data.Username})).Exists()
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

		if err = tx.WithContext(ctx).Create(&user).Error; err != nil {
			return fmt.Errorf("create user: %w", err)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return new(RegisterUser).fromUser(&user), nil
}

func (s *service) FindUserByCredential(ctx context.Context, data *LoginRequest) (*LoginUser, error) {
	user := new(model.User)
	if err := s.g.WithContext(ctx).Where(&model.User{Mobile: data.Mobile}).
		Or(&model.User{Username: data.Username}).
		Take(user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	if user.Disabled {
		return nil, ErrUserDisabled
	}

	if err := s.h.VerifyPassword(data.Password, user.Password); err != nil {
		return nil, ErrIncorrectCredentials
	}

	return new(LoginUser).fromUser(user), nil
}

func (u *LoginUser) fromUser(user *model.User) *LoginUser {
	u.ID = user.ID
	u.Name = user.Name
	u.Username = user.Username
	u.Mobile = user.Mobile
	return u
}

func (u *RegisterUser) fromUser(user *model.User) *RegisterUser {
	u.ID = user.ID
	u.Name = user.Name
	u.Username = user.Username
	u.Mobile = user.Mobile
	return u
}

func NewService(g *gorm.DB, h PasswordHash) Service {
	return &service{g: g, h: h}
}

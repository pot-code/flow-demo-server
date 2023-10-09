package auth

import (
	"context"
	"errors"
	"fmt"
	"gobit-demo/infra/event"
	"gobit-demo/infra/orm"
	"gobit-demo/model"
	"gobit-demo/services/auth/rbac"
	"gobit-demo/services/auth/session"
	"gobit-demo/services/auth/token"
	"time"

	"gorm.io/gorm"
)

var (
	ErrUserNotFound         = errors.New("用户不存在")
	ErrUserDisabled         = errors.New("用户已禁用")
	ErrDuplicatedUser       = errors.New("用户已存在")
	ErrIncorrectCredentials = errors.New("用户名或密码错误")
)

type Service interface {
	CreateUser(ctx context.Context, data *CreateUserDto) (*model.User, error)
	Login(ctx context.Context, data *LoginRequestDto) (string, error)
	GetUserPermissions(ctx context.Context, id model.ID) ([]string, error)
	GetUserRoles(ctx context.Context, id model.ID) ([]string, error)
}

type service struct {
	g  *gorm.DB
	h  PasswordHash
	sm session.SessionManager
	r  rbac.RoleService
	ts token.Service
	eb event.EventBus
}

func (s *service) CreateUser(ctx context.Context, data *CreateUserDto) (*model.User, error) {
	user := &model.User{
		Username: data.Username,
		Mobile:   data.Mobile,
		Name:     data.Name,
	}
	err := s.g.Transaction(func(tx *gorm.DB) error {
		exists, err := orm.Exists(
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

		role, err := s.r.GetRoleByName(ctx, "user")
		if err != nil {
			return fmt.Errorf("get role by name %s: %w", "user", err)
		}
		user.Roles = append(user.Roles, role)

		if err = tx.WithContext(ctx).Create(user).Error; err != nil {
			return fmt.Errorf("create user: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	s.eb.Publish(&UserCreatedEvent{
		UserID:    user.ID,
		Username:  user.Username,
		Mobile:    user.Mobile,
		Timestamp: time.Now().UnixMilli(),
	})

	return user, nil
}

func (s *service) Login(ctx context.Context, data *LoginRequestDto) (string, error) {
	m := new(model.User)
	err := s.g.WithContext(ctx).Where(&model.User{Mobile: data.Mobile}).
		Or(&model.User{Username: data.Username}).Take(m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "", ErrIncorrectCredentials
	}
	if err != nil {
		return "", fmt.Errorf("find user: %w", err)
	}

	if m.Disabled {
		return "", ErrUserDisabled
	}

	if err := s.h.VerifyPassword(data.Password, m.Password); err != nil {
		return "", ErrIncorrectCredentials
	}

	p, err := s.GetUserPermissions(ctx, m.ID)
	if err != nil {
		return "", fmt.Errorf("get user permissions: %w", err)
	}

	r, err := s.GetUserRoles(ctx, m.ID)
	if err != nil {
		return "", fmt.Errorf("get user roles: %w", err)
	}

	session, err := s.sm.NewSession(ctx, m.ID, m.Username, p, r)
	if err != nil {
		return "", fmt.Errorf("create session: %w", err)
	}

	token, err := s.ts.GenerateToken(&token.TokenData{SessionID: session.SessionID})
	if err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}
	return token, nil
}

func (s *service) GetUserPermissions(ctx context.Context, id model.ID) ([]string, error) {
	var permissions []string
	if err := s.g.WithContext(ctx).Model(&model.Permission{}).
		Distinct("permissions.name").
		Joins("INNER JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Joins("INNER JOIN user_roles ON user_roles.role_id = role_permissions.role_id").
		Where("user_roles.user_id = ?", id).
		Pluck("permissions.name", &permissions).Error; err != nil {
		return nil, fmt.Errorf("get user permissions: %w", err)
	}
	return permissions, nil
}

func (s *service) GetUserRoles(ctx context.Context, id model.ID) ([]string, error) {
	var roles []string
	if err := s.g.WithContext(ctx).Model(&model.Role{}).
		Joins("INNER JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", id).
		Pluck("roles.name", &roles).Error; err != nil {
		return nil, fmt.Errorf("get user roles: %w", err)
	}
	return roles, nil
}

func NewService(
	g *gorm.DB,
	eb event.EventBus,
	sm session.SessionManager,
	r rbac.RoleService,
	ts token.Service,
) Service {
	return &service{g: g, h: NewPasswordHash(), eb: eb, sm: sm, r: r, ts: ts}
}

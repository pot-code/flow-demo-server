package auth

import (
	"context"
	"fmt"
	"gobit-demo/internal/orm"
	"gobit-demo/model"

	"github.com/samber/lo"
	"gorm.io/gorm"
)

type UnAuthorizedError struct {
	UserID   uint   `json:"user_id,omitempty"`
	Username string `json:"username,omitempty"`
	Action   string `json:"action,omitempty"`
}

func (e UnAuthorizedError) Error() string {
	return fmt.Sprintf("no permission: username=%s, permission=%s", e.Username, e.Action)
}

type RBAC interface {
	CheckPermission(ctx context.Context, permission string) error
	CheckRole(ctx context.Context, role string) error
	GetRoles(ctx context.Context) ([]string, error)
	ListPermissions(ctx context.Context) ([]string, error)
	IsAdmin(ctx context.Context) (bool, error)
}

func NewRBAC(g *gorm.DB) RBAC {
	return &rbac{g: g}
}

type rbac struct {
	g *gorm.DB
}

func (s *rbac) IsAdmin(ctx context.Context) (bool, error) {
	roles, err := s.GetRoles(ctx)
	if err != nil {
		return false, fmt.Errorf("get roles: %w", err)
	}
	return lo.Contains(roles, "admin"), nil
}

func (s *rbac) CheckRole(ctx context.Context, role string) error {
	roles, err := s.GetRoles(ctx)
	if err != nil {
		return err
	}

	if lo.Contains(roles, role) {
		return nil
	}

	u := s.getLoginUser(ctx)
	return &UnAuthorizedError{
		UserID:   u.ID,
		Username: u.Username,
	}
}

func (s *rbac) CheckPermission(ctx context.Context, permission string) error {
	ok, err := s.IsAdmin(ctx)
	if err != nil {
		return fmt.Errorf("check admin: %w", err)
	}
	if ok {
		return nil
	}

	u := s.getLoginUser(ctx)
	ok, err = orm.NewGormWrapper(s.g.WithContext(ctx).Model(&model.User{}).
		Joins("INNER JOIN user_roles ur ON users.id = ur.user_id").
		Joins("INNER JOIN role_permissions rp ON ur.role_id = rp.role_id").
		Joins("INNER JOIN permissions p ON rp.permission_id = p.id").
		Where("users.id = ? AND p.name = ?", u.ID, permission)).Exists()
	if err != nil {
		return fmt.Errorf("get permission roles: %w", err)
	}
	if ok {
		return nil
	}
	return &UnAuthorizedError{
		UserID:   u.ID,
		Username: u.Username,
		Action:   permission,
	}
}

func (s *rbac) GetRoles(ctx context.Context) ([]string, error) {
	var roles []string

	u := s.getLoginUser(ctx)
	err := s.g.WithContext(ctx).Model(&model.User{}).
		Select("roles.name").
		Joins("INNER JOIN user_roles ON user_roles.user_id = users.id").
		Joins("INNER JOIN roles ON roles.id = user_roles.role_id").
		Where("users.id = ?", u.ID).
		Scan(&roles).Error
	return roles, err
}

func (s *rbac) ListPermissions(ctx context.Context) ([]string, error) {
	panic("unimplemented")
}

func (s *rbac) getLoginUser(ctx context.Context) *LoginUser {
	u, ok := new(LoginUser).FromContext(ctx)
	if !ok {
		panic(fmt.Errorf("no login user attached in context"))
	}
	return u
}

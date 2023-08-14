package auth

import (
	"context"
	"fmt"
	"gobit-demo/features/audit"
	"gobit-demo/model"

	"github.com/samber/lo"
	"gorm.io/gorm"
)

type UnAuthorizedError struct {
	UserID     uint   `json:"user_id,omitempty"`
	Username   string `json:"username,omitempty"`
	Permission string `json:"permission,omitempty"`
}

func (e UnAuthorizedError) Error() string {
	return fmt.Sprintf("no permission: username=%s, permission=%s", e.Username, e.Permission)
}

type RBAC interface {
	CheckPermission(ctx context.Context, permission string) error
	CheckRole(ctx context.Context, role string) error
	GetRoles(ctx context.Context) ([]string, error)
	GetPermissions(ctx context.Context) ([]string, error)
	IsAdmin(ctx context.Context) (bool, error)
}

func NewRBAC(g *gorm.DB, as audit.Service) RBAC {
	return &rbac{g: g, as: as}
}

type rbac struct {
	g  *gorm.DB
	as audit.Service
}

func (s *rbac) IsAdmin(ctx context.Context) (bool, error) {
	roles, err := s.GetRoles(ctx)
	if err != nil {
		return false, fmt.Errorf("get roles: %w", err)
	}
	return lo.Contains(roles, "admin"), nil
}

func (s *rbac) CheckRole(ctx context.Context, role string) error {
	u, ok := new(LoginUser).FromContext(ctx)
	if !ok {
		panic(fmt.Errorf("no login user attached in context"))
	}

	roles, err := s.GetRoles(ctx)
	if err != nil {
		return err
	}

	if lo.Contains(roles, role) {
		return nil
	}
	return &UnAuthorizedError{
		UserID:   u.ID,
		Username: u.Username,
	}
}

func (s *rbac) CheckPermission(ctx context.Context, permission string) error {
	u, ok := new(LoginUser).FromContext(ctx)
	if !ok {
		panic(fmt.Errorf("no login user attached in context"))
	}

	ok, err := s.IsAdmin(ctx)
	if err != nil {
		return fmt.Errorf("check admin: %w", err)
	}
	if ok {
		return nil
	}

	var allow []string
	if err := s.g.WithContext(ctx).Model(&model.Permission{}).
		Select("roles.name").
		Joins("INNER JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Joins("INNER JOIN roles ON roles.id = role_permissions.role_id").
		Where(&model.Permission{Name: permission}).
		Scan(&allow).Error; err != nil {
		return fmt.Errorf("get permission roles: %w", err)
	}

	roles, err := s.GetRoles(ctx)
	if err != nil {
		return err
	}

	for _, role := range roles {
		if lo.Contains(allow, role) {
			return nil
		}
	}

	re := &UnAuthorizedError{
		UserID:     u.ID,
		Username:   u.Username,
		Permission: permission,
	}
	if err := s.as.NewAuditLog().Subject(u.Username).Action("访问受限").Payload(re).
		Commit(ctx); err != nil {
		return fmt.Errorf("commit audit log: %w", err)
	}
	return re
}

func (s *rbac) GetRoles(ctx context.Context) ([]string, error) {
	u, ok := new(LoginUser).FromContext(ctx)
	if !ok {
		panic(fmt.Errorf("no login user attached in context"))
	}

	var roles []string
	err := s.g.WithContext(ctx).Model(&model.User{}).
		Select("roles.name").
		Joins("INNER JOIN user_roles ON user_roles.user_id = users.id").
		Joins("INNER JOIN roles ON roles.id = user_roles.role_id").
		Where("users.id = ?", u.ID).
		Scan(&roles).Error
	return roles, err
}

func (s *rbac) GetPermissions(ctx context.Context) ([]string, error) {
	panic("unimplemented")
}

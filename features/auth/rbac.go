package auth

import (
	"context"
	"fmt"
	"gobit-demo/model"

	"github.com/samber/lo"
	"gorm.io/gorm"
)

type UnAuthorizedError struct {
	UserID   uint
	Username string
	Obj      string
	Act      string
}

func (e UnAuthorizedError) Error() string {
	return fmt.Sprintf("no permission: username=%s, obj=%s, act=%s", e.Username, e.Obj, e.Act)
}

type RBAC interface {
	CheckPermission(ctx context.Context, obj, act string) error
	CheckRole(ctx context.Context, role string) error
	GetRoles(ctx context.Context) ([]string, error)
	GetPermissions(ctx context.Context) ([]string, error)
}

func NewRBAC(g *gorm.DB) RBAC {
	return &rbac{g: g}
}

type rbac struct {
	g *gorm.DB
}

func (*rbac) CheckRole(ctx context.Context, role string) error {
	panic("unimplemented")
}

func (s *rbac) CheckPermission(ctx context.Context, obj string, act string) error {
	u, ok := new(LoginUser).FromContext(ctx)
	if !ok {
		panic(fmt.Errorf("no login user attached in context"))
	}

	var allow []string
	if err := s.g.WithContext(ctx).Model(&model.Permission{}).
		Select("roles.name").
		Joins("INNER JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Joins("INNER JOIN roles ON roles.id = role_permissions.role_id").
		Where(&model.Permission{Object: obj, Action: act}).
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
	return &UnAuthorizedError{
		UserID:   u.ID,
		Username: u.Username,
		Obj:      obj,
		Act:      act,
	}
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

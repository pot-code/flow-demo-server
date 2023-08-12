package rbac

import (
	"context"
	"fmt"
	"gobit-demo/features/auth"
	"gobit-demo/model"

	"github.com/samber/lo"
	"gorm.io/gorm"
)

type Service interface {
	HasPermission(ctx context.Context, obj, act string) (bool, error)
	HasRole(ctx context.Context, role string) error
	GetRoles(ctx context.Context) ([]string, error)
	GetPermissions(ctx context.Context) ([]string, error)
}

func NewService(g *gorm.DB) Service {
	return &service{g: g}
}

type service struct {
	g *gorm.DB
}

func (*service) HasRole(ctx context.Context, role string) error {
	panic("unimplemented")
}

func (s *service) HasPermission(ctx context.Context, obj string, act string) (bool, error) {
	var allow []string
	if err := s.g.WithContext(ctx).Model(&model.Permission{}).
		Select("roles.name").
		Joins("INNER JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Joins("INNER JOIN roles ON roles.id = role_permissions.role_id").
		Where(&model.Permission{Object: obj, Action: act}).
		Scan(&allow).Error; err != nil {
		return false, fmt.Errorf("get permission roles: %w", err)
	}

	roles, err := s.GetRoles(ctx)
	if err != nil {
		return false, err
	}

	for _, role := range roles {
		if lo.Contains(allow, role) {
			return true, nil
		}
	}
	return false, nil
}

func (s *service) GetRoles(ctx context.Context) ([]string, error) {
	u, ok := new(auth.LoginUser).FromContext(ctx)
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

func (s *service) GetPermissions(ctx context.Context) ([]string, error) {
	panic("unimplemented")
}

type PermissionService interface {
	AddPermission(ctx context.Context, role, obj, act string) error
	UpdatePermission(ctx context.Context, role, obj, act string) error
	DeletePermission(ctx context.Context, role, obj, act string) error
}

func NewPermissionService(g *gorm.DB) PermissionService {
	return &permissionService{g: g}
}

type permissionService struct {
	g *gorm.DB
}

// AddPermission implements PolicyService.
func (*permissionService) AddPermission(ctx context.Context, role string, obj string, act string) error {
	panic("unimplemented")
}

// DeletePermission implements PolicyService.
func (*permissionService) DeletePermission(ctx context.Context, role string, obj string, act string) error {
	panic("unimplemented")
}

// UpdatePermission implements PolicyService.
func (*permissionService) UpdatePermission(ctx context.Context, role string, obj string, act string) error {
	panic("unimplemented")
}

type RoleService interface {
	AddRole(ctx context.Context, role string) error
	UpdateRole(ctx context.Context, role string) error
	DeleteRole(ctx context.Context, role string) error
}

func NewRoleService(g *gorm.DB) RoleService {
	return &roleService{g: g}
}

type roleService struct {
	g *gorm.DB
}

// AddRole implements RoleService.
func (*roleService) AddRole(ctx context.Context, role string) error {
	panic("unimplemented")
}

// DeleteRole implements RoleService.
func (*roleService) DeleteRole(ctx context.Context, role string) error {
	panic("unimplemented")
}

// UpdateRole implements RoleService.
func (*roleService) UpdateRole(ctx context.Context, role string) error {
	panic("unimplemented")
}

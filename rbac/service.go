package role

import (
	"context"

	"gorm.io/gorm"
)

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

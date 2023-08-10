package rbac

import (
	"context"
	"fmt"
	"gobit-demo/features/auth"
	"strconv"

	"github.com/casbin/casbin/v2"
	"gorm.io/gorm"
)

type Service interface {
	HasPermission(ctx context.Context, obj, act string) error
}

func NewService(g *gorm.DB) Service {
	return &service{e: newCasbinEnforcer(g)}
}

type service struct {
	e *casbin.Enforcer
}

func (s *service) HasPermission(ctx context.Context, obj string, act string) error {
	u, ok := new(auth.LoginUser).FromContext(ctx)
	if !ok {
		panic(fmt.Errorf("no login user attached in context"))
	}

	ok, err := s.e.Enforce(strconv.Itoa(int(u.Id)), obj, act)
	if err != nil {
		return fmt.Errorf("check permission: %w", err)
	}
	if !ok {
		return &NoPermissionError{
			UserID:   u.Id,
			Username: u.Username,
			Obj:      obj,
			Act:      act,
		}
	}
	return nil
}

type PolicyService interface {
	AddPolicy(ctx context.Context, role, obj, act string) error
	UpdatePolicy(ctx context.Context, role, obj, act string) error
	DeletePolicy(ctx context.Context, role, obj, act string) error
}

func NewPolicyService(g *gorm.DB) PolicyService {
	return &policyService{g: g}
}

type policyService struct {
	g *gorm.DB
}

// AddPolicy implements PolicyService.
func (*policyService) AddPolicy(ctx context.Context, role string, obj string, act string) error {
	panic("unimplemented")
}

// DeletePolicy implements PolicyService.
func (*policyService) DeletePolicy(ctx context.Context, role string, obj string, act string) error {
	panic("unimplemented")
}

// UpdatePolicy implements PolicyService.
func (*policyService) UpdatePolicy(ctx context.Context, role string, obj string, act string) error {
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

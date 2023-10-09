package auth

import (
	"context"
	"fmt"
	"gobit-demo/services/auth/session"

	"github.com/samber/lo"
	"gorm.io/gorm"
)

type RBAC interface {
	CheckPermission(ctx context.Context, permission string) error
	CheckRole(ctx context.Context, role string) error
	IsAdmin(ctx context.Context) (bool, error)
}

type rbac struct {
	g *gorm.DB
}

func (r *rbac) IsAdmin(ctx context.Context) (bool, error) {
	s := session.GetSessionFromContext(ctx)
	if lo.Contains(s.UserRoles, "admin") {
		return true, nil
	}
	return false, nil
}

func (r *rbac) CheckRole(ctx context.Context, role string) error {
	ok, err := r.IsAdmin(ctx)
	if err != nil {
		return fmt.Errorf("check admin: %w", err)
	}
	if ok {
		return nil
	}

	s := session.GetSessionFromContext(ctx)
	if lo.Contains(s.UserRoles, role) {
		return nil
	}
	return new(UnAuthorizedError)
}

func (r *rbac) CheckPermission(ctx context.Context, permission string) error {
	ok, err := r.IsAdmin(ctx)
	if err != nil {
		return fmt.Errorf("check admin: %w", err)
	}
	if ok {
		return nil
	}

	s := session.GetSessionFromContext(ctx)
	if lo.Contains(s.UserPermissions, permission) {
		return nil
	}
	return new(UnAuthorizedError)
}

func NewRBAC(g *gorm.DB) RBAC {
	return &rbac{g: g}
}

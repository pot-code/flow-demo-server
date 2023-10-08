package auth

import (
	"context"
	"fmt"

	"github.com/samber/lo"
	"gorm.io/gorm"
)

type RBAC interface {
	CheckPermission(ctx context.Context, permission string) error
	CheckRole(ctx context.Context, role string) error
	IsAdmin(ctx context.Context) (bool, error)
}

type rbac struct {
	g  *gorm.DB
	sm SessionManager
}

func (r *rbac) IsAdmin(ctx context.Context) (bool, error) {
	s := r.sm.GetSessionFromContext(ctx)
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

	s := r.sm.GetSessionFromContext(ctx)
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

	s := r.sm.GetSessionFromContext(ctx)
	if lo.Contains(s.UserPermissions, permission) {
		return nil
	}
	return new(UnAuthorizedError)
}

func NewRBAC(g *gorm.DB, sm SessionManager) RBAC {
	return &rbac{g: g, sm: sm}
}

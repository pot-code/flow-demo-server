package auth

import (
	"context"
	"fmt"
	"gobit-demo/model"

	"github.com/samber/lo"
	"gorm.io/gorm"
)

type UnAuthorizedError struct {
	UserID model.ID `json:"user_id,omitempty"`
	Action string   `json:"action,omitempty"`
}

func (e UnAuthorizedError) Error() string {
	return fmt.Sprintf("no permission: user_id=%v, permission=%s", e.UserID, e.Action)
}

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
	return &UnAuthorizedError{UserID: s.UserID}
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
	return &UnAuthorizedError{
		UserID: s.UserID,
		Action: permission,
	}
}

func NewRBAC(g *gorm.DB, sm SessionManager) RBAC {
	return &rbac{g: g, sm: sm}
}

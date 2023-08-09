package perm

import (
	"context"
)

type PermService interface {
	HasPerm(ctx context.Context, obj, act string) bool
	AddPerm(ctx context.Context, role, obj, act string) bool
	DeletePerm(ctx context.Context, role, obj, act string) bool
}

type RoleService interface {
	HasRole(ctx context.Context, role string) bool
}

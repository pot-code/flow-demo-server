package perm

import "context"

type RoleService interface {
	HasRole(ctx context.Context, role string) bool
}

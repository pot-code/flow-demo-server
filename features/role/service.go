package role

import (
	"context"

	"gorm.io/gorm"
)

type Service interface {
	AddRole(ctx context.Context, role string) error
	UpdateRole(ctx context.Context, role string) error
	DeleteRole(ctx context.Context, role string) error
}

func NewService(g *gorm.DB) Service {
	return &service{g: g}
}

type service struct {
	g *gorm.DB
}

// AddRole implements RoleService.
func (s *service) AddRole(ctx context.Context, role string) error {
	panic("unimplemented")
}

// DeleteRole implements RoleService.
func (s *service) DeleteRole(ctx context.Context, role string) error {
	panic("unimplemented")
}

// UpdateRole implements RoleService.
func (s *service) UpdateRole(ctx context.Context, role string) error {
	panic("unimplemented")
}

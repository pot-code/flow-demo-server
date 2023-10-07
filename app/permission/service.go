package role

import (
	"context"

	"gorm.io/gorm"
)

type Service interface {
	AddPermission(ctx context.Context, role, obj, act string) error
	UpdatePermission(ctx context.Context, role, obj, act string) error
	DeletePermission(ctx context.Context, role, obj, act string) error
}

func NewService(g *gorm.DB) Service {
	return &service{g: g}
}

type service struct {
	g *gorm.DB
}

// AddPermission implements PolicyService.
func (s *service) AddPermission(ctx context.Context, role string, obj string, act string) error {
	panic("unimplemented")
}

// DeletePermission implements PolicyService.
func (s *service) DeletePermission(ctx context.Context, role string, obj string, act string) error {
	panic("unimplemented")
}

// UpdatePermission implements PolicyService.
func (s *service) UpdatePermission(ctx context.Context, role string, obj string, act string) error {
	panic("unimplemented")
}

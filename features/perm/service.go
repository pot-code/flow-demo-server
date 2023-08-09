package perm

import (
	"context"

	"gorm.io/gorm"
)

type Service interface {
	HasPerm(ctx context.Context, obj, act string) error
	AddPerm(ctx context.Context, role, obj, act string) error
	DeletePerm(ctx context.Context, role, obj, act string) error
}

func NewService(db *gorm.DB) Service {
	return nil
}

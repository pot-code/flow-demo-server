package rbac

import (
	"context"
	"fmt"
	"gobit-demo/model"

	"gorm.io/gorm"
)

type RoleService interface {
	GetRoleByName(ctx context.Context, name string) (*model.Role, error)
}

type roleService struct {
	g *gorm.DB
}

func (r *rbac) GetRoleByName(ctx context.Context, name string) (*model.Role, error) {
	rm := new(model.Role)
	if err := r.g.WithContext(ctx).Where(&model.Role{Name: name}).Take(rm).Error; err != nil {
		return nil, fmt.Errorf("get role: %w", err)
	}
	return rm, nil
}

func NewRoleService(g *gorm.DB) *roleService {
	return &roleService{g: g}
}

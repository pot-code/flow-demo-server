package flow

import (
	"context"
	"fmt"
	"gobit-demo/features/auth"
	"gobit-demo/internal/orm"
	"gobit-demo/model"

	"gorm.io/gorm"
)

type PermissionService interface {
	CanViewFlow(ctx context.Context, fid model.UUID) error
	CanUpdateFlow(ctx context.Context, fid model.UUID) error
	CanDeleteFlow(ctx context.Context, fid model.UUID) error
}

type permission struct {
	g  *gorm.DB
	sm auth.SessionManager
}

// CanDeleteFlow implements PermissionService.
func (p *permission) CanDeleteFlow(ctx context.Context, fid model.UUID) error {
	return p.CanViewFlow(ctx, fid)
}

func (p *permission) CanUpdateFlow(ctx context.Context, fid model.UUID) error {
	return p.CanViewFlow(ctx, fid)
}

func (p *permission) CanViewFlow(ctx context.Context, fid model.UUID) error {
	s := p.sm.GetSessionFromContext(ctx)
	ok, err := new(orm.GormUtil).Exists(p.g.WithContext(ctx).Model(&model.Flow{}).Where("id = ? AND owner_id = ?", fid, s.UserID))
	if err != nil {
		return fmt.Errorf("check flow exists by id: %w", err)
	}
	if !ok {
		return &auth.UnAuthorizedError{
			UserID: s.UserID,
			Action: fmt.Sprintf("view flow %v", fid),
		}
	}
	return nil
}

func NewPermissionService(g *gorm.DB, sm auth.SessionManager) PermissionService {
	return &permission{g: g, sm: sm}
}

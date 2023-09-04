package flow

import (
	"context"
	"fmt"
	"gobit-demo/auth"
	"gobit-demo/internal/orm"
	"gobit-demo/model"

	"gorm.io/gorm"
)

type ABAC interface {
	CanViewFlow(ctx context.Context, fid model.UUID) error
	CanUpdateFlow(ctx context.Context, fid model.UUID) error
	CanDeleteFlow(ctx context.Context, fid model.UUID) error
}

type abac struct {
	g  *gorm.DB
	sm auth.SessionManager
}

// CanDeleteFlow implements PermissionService.
func (p *abac) CanDeleteFlow(ctx context.Context, fid model.UUID) error {
	return p.CanViewFlow(ctx, fid)
}

func (p *abac) CanUpdateFlow(ctx context.Context, fid model.UUID) error {
	return p.CanViewFlow(ctx, fid)
}

func (p *abac) CanViewFlow(ctx context.Context, fid model.UUID) error {
	s := p.sm.GetSessionFromContext(ctx)
	ok, err := orm.Exists(p.g.WithContext(ctx).Model(&model.Flow{}).Where("id = ? AND owner_id = ?", fid, s.UserID))
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

func NewABAC(g *gorm.DB, sm auth.SessionManager) ABAC {
	return &abac{g: g, sm: sm}
}

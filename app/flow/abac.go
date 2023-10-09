package flow

import (
	"context"
	"fmt"
	"gobit-demo/infra/orm"
	"gobit-demo/model"
	"gobit-demo/services/auth/rbac"
	"gobit-demo/services/auth/session"

	"gorm.io/gorm"
)

type ABAC struct {
	g *gorm.DB
}

func (p *ABAC) CanDelete(ctx context.Context, id model.ID) error {
	ok, err := p.isOwner(ctx, id)
	if err != nil {
		return err
	}
	if !ok {
		return rbac.ErrUnAuthorized
	}
	return nil
}

func (p *ABAC) isOwner(ctx context.Context, id model.ID) (bool, error) {
	s := session.GetSessionFromContext(ctx)
	ok, err := orm.Exists(p.g.WithContext(ctx).Model(&model.Flow{}).Where("id = ? AND owner_id = ?", id, s.UserID))
	if err != nil {
		return false, fmt.Errorf("check flow exists by id: %w", err)
	}
	return ok, err
}

func NewABAC(g *gorm.DB) *ABAC {
	return &ABAC{g: g}
}

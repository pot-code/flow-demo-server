package flow

import (
	"context"
	"fmt"
	"gobit-demo/infra/orm"
	"gobit-demo/model"
	"gobit-demo/services/auth"

	"gorm.io/gorm"
)

type ABAC interface {
	CanView(ctx context.Context, id model.ID) error
	CanUpdate(ctx context.Context, id model.ID) error
	CanDelete(ctx context.Context, id model.ID) error
}

type abac struct {
	g  *gorm.DB
	sm auth.SessionManager
}

func (p *abac) CanDelete(ctx context.Context, id model.ID) error {
	ok, err := p.isOwner(ctx, id)
	if err != nil {
		return err
	}
	if !ok {
		return new(auth.UnAuthorizedError)
	}
	return nil

}

func (p *abac) CanUpdate(ctx context.Context, id model.ID) error {
	ok, err := p.isOwner(ctx, id)
	if err != nil {
		return err
	}
	if !ok {
		return new(auth.UnAuthorizedError)
	}
	return nil

}

func (p *abac) CanView(ctx context.Context, id model.ID) error {
	ok, err := p.isOwner(ctx, id)
	if err != nil {
		return err
	}
	if !ok {
		return new(auth.UnAuthorizedError)
	}
	return nil
}

func (p *abac) isOwner(ctx context.Context, id model.ID) (bool, error) {
	s := p.sm.GetSessionFromContext(ctx)
	ok, err := orm.Exists(p.g.WithContext(ctx).Model(&model.Flow{}).Where("id = ? AND owner_id = ?", id, s.UserID))
	if err != nil {
		return false, fmt.Errorf("check flow exists by id: %w", err)
	}
	return ok, err
}

func NewPermission(g *gorm.DB, sm auth.SessionManager) *abac {
	return &abac{g: g, sm: sm}
}

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
	CanViewFlowByID(ctx context.Context, fid model.UUID) error
}

type permission struct {
	g  *gorm.DB
	sm auth.SessionManager
}

func (p *permission) CanViewFlowByID(ctx context.Context, fid model.UUID) error {
	s := p.sm.GetSessionFromContext(ctx)
	ok, err := orm.NewGormWrapper(p.g.WithContext(ctx).Model(&model.Flow{}).
		Where("id = ? AND owner_id = ?", fid, s.UserID)).Exists()
	if err != nil {
		return fmt.Errorf("check flow exists by id: %w", err)
	}
	if !ok {
		return &auth.UnAuthorizedError{
			UserID: s.UserID,
			Action: fmt.Sprintf("view flow %s", fid),
		}
	}
	return nil
}

func NewPermissionService(g *gorm.DB, sm auth.SessionManager) PermissionService {
	return &permission{g: g, sm: sm}
}

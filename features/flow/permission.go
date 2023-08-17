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
	CanViewFlowByID(ctx context.Context, fid uint) error
}

type permission struct {
	g *gorm.DB
}

func (p *permission) CanViewFlowByID(ctx context.Context, fid uint) error {
	s, _ := new(auth.Session).FromContext(ctx)
	ok, err := orm.NewGormWrapper(p.g.WithContext(ctx).Model(&model.Flow{}).Where("id = ? AND owner_id = ?", fid, s.UserID)).Exists()
	if err != nil {
		return fmt.Errorf("check flow exists by id: %w", err)
	}
	if !ok {
		return &auth.UnAuthorizedError{
			UserID: s.UserID,
			Action: fmt.Sprintf("view flow %d", fid),
		}
	}
	return nil
}

func NewPermissionService(g *gorm.DB) PermissionService {
	return &permission{g: g}
}

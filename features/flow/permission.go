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
	u := p.getLoginUser(ctx)
	ok, err := orm.NewGormWrapper(p.g.WithContext(ctx).Model(&model.Flow{}).Where("id = ? AND owner_id = ?", fid, u.ID)).Exists()
	if err != nil {
		return fmt.Errorf("check flow exists by id: %w", err)
	}
	if !ok {
		return &auth.UnAuthorizedError{
			UserID:   u.ID,
			Username: u.Username,
			Action:   fmt.Sprintf("view flow %d", fid),
		}
	}
	return nil
}

func (p *permission) getLoginUser(ctx context.Context) *auth.LoginUser {
	u, ok := new(auth.LoginUser).FromContext(ctx)
	if !ok {
		panic(fmt.Errorf("no login user attached in context"))
	}
	return u
}

func NewPermissionService(g *gorm.DB) PermissionService {
	return &permission{g: g}
}

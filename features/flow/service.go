package flow

import (
	"context"
	"errors"
	"fmt"
	"gobit-demo/features/auth"
	"gobit-demo/internal/pagination"
	"gobit-demo/model"
	"gobit-demo/util"

	"gorm.io/gorm"
)

var (
	ErrDuplicatedFlow = errors.New("流程已存在")
)

type Service interface {
	GetFlowByID(ctx context.Context, fid model.UUID) (*model.Flow, error)
	ListFlow(ctx context.Context, p *pagination.Pagination) ([]*model.Flow, int, error)
	CreateFlow(ctx context.Context, req *CreateFlowRequest) error
	UpdateFlow(ctx context.Context, req *UpdateFlowRequest) error
	DeleteFlow(ctx context.Context, fid model.UUID) error
}

type service struct {
	g  *gorm.DB
	a  ABAC
	sm auth.SessionManager
}

// DeleteFlow implements Service.
func (s *service) DeleteFlow(ctx context.Context, fid model.UUID) error {
	if err := s.a.CanDeleteFlow(ctx, fid); err != nil {
		return err
	}

	return s.g.WithContext(ctx).Model(&model.Flow{}).Where("id = ?", fid).Delete(&model.Flow{}).Error
}

func (s *service) GetFlowByID(ctx context.Context, fid model.UUID) (*model.Flow, error) {
	if err := s.a.CanViewFlow(ctx, fid); err != nil {
		return nil, err
	}

	m := new(model.Flow)
	if err := s.g.WithContext(ctx).Model(&model.Flow{}).Where("id = ?", fid).Take(m).Error; err != nil {
		return nil, fmt.Errorf("get flow by id: %w", err)
	}
	return m, nil
}

func (s *service) CreateFlow(ctx context.Context, req *CreateFlowRequest) error {
	session := s.sm.GetSessionFromContext(ctx)
	m := &model.Flow{
		Name:        req.Name,
		Nodes:       req.Nodes,
		Edges:       req.Edges,
		Description: req.Description,
		OwnerID:     &session.UserID,
	}
	if err := s.g.WithContext(ctx).Create(m).Error; err != nil {
		return fmt.Errorf("create flow: %w", err)
	}
	return nil
}

func (s *service) UpdateFlow(ctx context.Context, req *UpdateFlowRequest) error {
	if err := s.a.CanUpdateFlow(ctx, req.ID); err != nil {
		return err
	}

	m := &model.Flow{
		ID:          req.ID,
		Name:        req.Name,
		Nodes:       req.Nodes,
		Edges:       req.Edges,
		Description: req.Description,
	}
	if err := s.g.WithContext(ctx).Model(&model.Flow{}).Where(&model.Flow{ID: req.ID}).Updates(m).Error; err != nil {
		return fmt.Errorf("update flow: %w", err)
	}
	return nil
}

func (s *service) ListFlow(ctx context.Context, p *pagination.Pagination) ([]*model.Flow, int, error) {
	var (
		flows []*model.Flow
		count int64
	)

	if err := s.g.WithContext(ctx).Model(&model.Flow{}).
		Scopes(new(util.GormUtil).Pagination(p)).
		Find(&flows).
		Count(&count).
		Error; err != nil {
		return nil, 0, fmt.Errorf("query flow list: %w", err)
	}
	return flows, int(count), nil
}

func NewService(g *gorm.DB, sm auth.SessionManager, p ABAC) Service {
	return &service{g: g, sm: sm, a: p}
}

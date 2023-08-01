package flow

import (
	"context"
	"errors"
	"fmt"
	"gobit-demo/internal/pagination"
	"gobit-demo/model"

	"gorm.io/gorm"
)

var (
	ErrDuplicatedFlow = errors.New("流程已存在")
)

type FlowService struct {
	g *gorm.DB
}

func NewFlowService(g *gorm.DB) *FlowService {
	return &FlowService{g: g}
}

func (s *FlowService) CreateFlow(ctx context.Context, req *CreateFlowRequest) error {
	var count int64
	if err := s.g.WithContext(ctx).Model(&model.Flow{}).
		Where(&model.Flow{Name: req.Name}).
		Count(&count).
		Error; err != nil {
		return fmt.Errorf("check duplicate flow: %w", err)
	}
	if count > 0 {
		return ErrDuplicatedFlow
	}

	save := &model.Flow{
		Name:        req.Name,
		Description: req.Description,
	}
	if err := s.g.Create(save).Error; err != nil {
		return fmt.Errorf("create flow: %w", err)
	}
	return nil
}

func (s *FlowService) ListFlow(ctx context.Context, p *pagination.Pagination) ([]*ListFlowDto, uint, error) {
	var (
		flows []*ListFlowDto
		count int64
	)

	if err := s.g.WithContext(ctx).Model(&model.Flow{}).
		Limit(p.PageSize).
		Offset((p.Page - 1) * p.PageSize).
		Find(&flows).
		Count(&count).
		Error; err != nil {
		return nil, 0, fmt.Errorf("query flow list: %w", err)
	}
	return flows, uint(count), nil
}

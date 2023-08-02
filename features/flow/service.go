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
	ErrDuplicatedFlow     = errors.New("流程已存在")
	ErrDuplicatedFlowNode = errors.New("流程结点已存在")
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
	if err := s.g.WithContext(ctx).Create(save).Error; err != nil {
		return fmt.Errorf("create flow: %w", err)
	}
	return nil
}

func (s *FlowService) ListFlow(ctx context.Context, p *pagination.Pagination) ([]*ListFlowResponse, int, error) {
	var (
		flows []*ListFlowResponse
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
	return flows, int(count), nil
}

func (s *FlowService) CreateFlowNode(ctx context.Context, req *CreateFlowNodeRequest) error {
	var nodes []*model.FlowNode
	if err := s.g.WithContext(ctx).Model(&model.FlowNode{}).
		Where(&model.FlowNode{FlowID: *req.FlowID}).
		Find(&nodes).
		Error; err != nil {
		return fmt.Errorf("query flow node: %w", err)
	}
	for _, node := range nodes {
		if node.Name == req.Name {
			return ErrDuplicatedFlowNode
		}
	}

	if err := s.g.WithContext(ctx).Create(&model.FlowNode{
		FlowID:      *req.FlowID,
		Name:        req.Name,
		Description: req.Description,
		PrevID:      req.PrevID,
		NextID:      req.NextID,
	}).Error; err != nil {
		return fmt.Errorf("create flow node: %w", err)
	}
	return nil
}

func (s *FlowService) ListFlowNodeByFlowID(ctx context.Context, flowID uint) ([]*ListFlowNodeResponse, error) {
	var nodes []*ListFlowNodeResponse
	if err := s.g.WithContext(ctx).Model(&model.FlowNode{}).
		Where(&model.FlowNode{FlowID: flowID}).
		Find(&nodes).
		Error; err != nil {
		return nil, fmt.Errorf("query flow node: %w", err)
	}
	return nodes, nil
}

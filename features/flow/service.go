package flow

import (
	"context"
	"errors"
	"fmt"
	"gobit-demo/internal/orm"
	"gobit-demo/internal/pagination"
	"gobit-demo/model"

	"gorm.io/gorm"
)

var (
	ErrDuplicatedFlow     = errors.New("流程已存在")
	ErrDuplicatedFlowNode = errors.New("流程结点已存在")
)

type Service interface {
	CreateFlow(ctx context.Context, data *CreateFlowRequest) error
	ListFlow(ctx context.Context, p *pagination.Pagination) ([]*ListFlowResponse, int, error)
	CreateFlowNode(ctx context.Context, data *CreateFlowNodeRequest) error
	ListFlowNodeByFlowID(ctx context.Context, flowID uint) ([]*ListFlowNodeResponse, error)
}

type service struct {
	g *gorm.DB
}

func NewService(g *gorm.DB) Service {
	return &service{g: g}
}

func (s *service) CreateFlow(ctx context.Context, req *CreateFlowRequest) error {
	return s.g.Transaction(func(tx *gorm.DB) error {
		exists, err := orm.NewGormWrapper(tx.WithContext(ctx).Model(&model.Flow{}).
			Where(&model.Flow{Name: req.Name})).Exists()
		if err != nil {
			return fmt.Errorf("check duplicate flow: %w", err)
		}
		if exists {
			return ErrDuplicatedFlow
		}

		save := &model.Flow{
			Name:        req.Name,
			Description: req.Description,
		}
		if err := tx.WithContext(ctx).Create(save).Error; err != nil {
			return fmt.Errorf("create flow: %w", err)
		}
		return nil
	})
}

func (s *service) ListFlow(ctx context.Context, p *pagination.Pagination) ([]*ListFlowResponse, int, error) {
	var (
		flows []*ListFlowResponse
		count int64
	)

	if err := orm.NewGormWrapper(s.g.WithContext(ctx).Model(&model.Flow{})).
		Paginate(p).
		Find(&flows).
		Count(&count).
		Error; err != nil {
		return nil, 0, fmt.Errorf("query flow list: %w", err)
	}
	return flows, int(count), nil
}

func (s *service) CreateFlowNode(ctx context.Context, req *CreateFlowNodeRequest) error {
	return s.g.Transaction(func(tx *gorm.DB) error {
		ok, err := orm.NewGormWrapper(tx.WithContext(ctx).Model(&model.FlowNode{}).
			Where(&model.FlowNode{FlowID: *req.FlowID})).Exists()
		if err != nil {
			return fmt.Errorf("check duplicate flow node: %w", err)
		}
		if ok {
			return ErrDuplicatedFlowNode
		}

		if err := tx.WithContext(ctx).Create(&model.FlowNode{
			FlowID:      *req.FlowID,
			Name:        req.Name,
			Description: req.Description,
			PrevID:      req.PrevID,
			NextID:      req.NextID,
		}).Error; err != nil {
			return fmt.Errorf("create flow node: %w", err)
		}
		return nil
	})
}

func (s *service) ListFlowNodeByFlowID(ctx context.Context, flowID uint) ([]*ListFlowNodeResponse, error) {
	var nodes []*ListFlowNodeResponse
	if err := s.g.WithContext(ctx).Model(&model.FlowNode{}).
		Where(&model.FlowNode{FlowID: flowID}).
		Find(&nodes).
		Error; err != nil {
		return nil, fmt.Errorf("query flow node: %w", err)
	}
	return nodes, nil
}

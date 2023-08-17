package flow

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gobit-demo/internal/orm"
	"gobit-demo/internal/pagination"
	"gobit-demo/model"

	"gorm.io/gorm"
)

var (
	ErrDuplicatedFlow = errors.New("流程已存在")
)

type Service interface {
	CreateFlow(ctx context.Context, req *CreateFlowRequest) error
	UpdateFlow(ctx context.Context, req *UpdateFlowRequest) error
	ListFlow(ctx context.Context, p *pagination.Pagination) ([]*ListFlowResponse, int, error)
	GetFlowByID(ctx context.Context, fid uint) (*FlowObjectResponse, error)
}

type service struct {
	g *gorm.DB
}

func NewService(g *gorm.DB) Service {
	return &service{g: g}
}

func (s *service) GetFlowByID(ctx context.Context, fid uint) (*FlowObjectResponse, error) {
	m := new(model.Flow)
	if err := s.g.WithContext(ctx).Model(&model.Flow{}).Where("id = ?", fid).Take(m).Error; err != nil {
		return nil, fmt.Errorf("get flow by id: %w", err)
	}

	o := new(FlowObjectResponse)
	if err := json.Unmarshal([]byte(m.Nodes), o); err != nil {
		return nil, fmt.Errorf("unmarshal nodes: %w", err)
	}
	if err := json.Unmarshal([]byte(m.Edges), o); err != nil {
		return nil, fmt.Errorf("unmarshal edges: %w", err)
	}
	return o, nil
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

		nodes, err := json.Marshal(req.Nodes)
		if err != nil {
			return fmt.Errorf("marshal nodes: %w", err)
		}
		edges, err := json.Marshal(req.Edges)
		if err != nil {
			return fmt.Errorf("marshal edges: %w", err)
		}
		save := &model.Flow{
			Name:        req.Name,
			Nodes:       string(nodes),
			Edges:       string(edges),
			Description: req.Description,
		}
		if err := tx.WithContext(ctx).Create(save).Error; err != nil {
			return fmt.Errorf("create flow: %w", err)
		}
		return nil
	})
}

func (s *service) UpdateFlow(ctx context.Context, req *UpdateFlowRequest) error {
	nodes, err := json.Marshal(req.Nodes)
	if err != nil {
		return fmt.Errorf("marshal nodes: %w", err)
	}
	edges, err := json.Marshal(req.Edges)
	if err != nil {
		return fmt.Errorf("marshal edges: %w", err)
	}
	save := &model.Flow{
		Model:       gorm.Model{ID: *req.ID},
		Name:        req.Name,
		Nodes:       string(nodes),
		Edges:       string(edges),
		Description: req.Description,
	}
	if err := s.g.WithContext(ctx).Model(&model.Flow{}).Save(save).Error; err != nil {
		return fmt.Errorf("update flow: %w", err)
	}
	return nil
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

package flow

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gobit-demo/features/auth"
	"gobit-demo/internal/orm"
	"gobit-demo/internal/pagination"
	"gobit-demo/model"

	"gorm.io/gorm"
)

var (
	ErrDuplicatedFlow = errors.New("流程已存在")
)

type Service interface {
	GetFlowByID(ctx context.Context, fid string) (*FlowObjectResponse, error)
	ListFlow(ctx context.Context, p *pagination.Pagination) ([]*ListFlowResponse, int, error)
	CreateFlow(ctx context.Context, req *CreateFlowRequest) error
	UpdateFlow(ctx context.Context, req *UpdateFlowRequest) error
}

type service struct {
	g  *gorm.DB
	sm auth.SessionManager
}

func NewService(g *gorm.DB, sm auth.SessionManager) Service {
	return &service{g: g, sm: sm}
}

func (s *service) GetFlowByID(ctx context.Context, fid string) (*FlowObjectResponse, error) {
	m := new(model.Flow)
	if err := s.g.WithContext(ctx).Model(&model.Flow{}).Where("id = ?", fid).Take(m).Error; err != nil {
		return nil, fmt.Errorf("get flow by id: %w", err)
	}

	o := new(FlowObjectResponse)
	if err := json.Unmarshal([]byte(m.Nodes), &o.Nodes); err != nil {
		return nil, fmt.Errorf("unmarshal nodes: %w", err)
	}
	if err := json.Unmarshal([]byte(m.Edges), &o.Edges); err != nil {
		return nil, fmt.Errorf("unmarshal edges: %w", err)
	}
	return o, nil
}

func (s *service) CreateFlow(ctx context.Context, req *CreateFlowRequest) error {
	session := s.sm.GetSessionFromContext(ctx)
	nodes, err := json.Marshal(req.Nodes)
	if err != nil {
		return fmt.Errorf("marshal nodes: %w", err)
	}
	edges, err := json.Marshal(req.Edges)
	if err != nil {
		return fmt.Errorf("marshal edges: %w", err)
	}
	m := &model.Flow{
		Name:        req.Name,
		Nodes:       string(nodes),
		Edges:       string(edges),
		Description: req.Description,
		OwnerID:     &session.UserID,
	}
	if err := s.g.WithContext(ctx).Create(m).Error; err != nil {
		return fmt.Errorf("create flow: %w", err)
	}
	return nil
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
	m := &model.Flow{
		ID:          req.ID,
		Name:        req.Name,
		Nodes:       string(nodes),
		Edges:       string(edges),
		Description: req.Description,
	}
	if err := s.g.WithContext(ctx).Model(&model.Flow{}).Where(&model.Flow{ID: req.ID}).Updates(m).Error; err != nil {
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

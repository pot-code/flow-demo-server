package flow

import (
	"context"
	"errors"
	"fmt"
	"gobit-demo/audit"
	"gobit-demo/auth"
	"gobit-demo/internal/event"
	"gobit-demo/internal/orm"
	"gobit-demo/internal/pagination"
	"gobit-demo/model"
	"time"

	"gorm.io/gorm"
)

var (
	ErrDuplicatedFlow = errors.New("流程已存在")
)

type Service interface {
	GetFlowByID(ctx context.Context, fid model.UUID) (*model.Flow, error)
	ListFlowByOwner(ctx context.Context, p *pagination.Pagination) ([]*model.Flow, int, error)
	CreateFlow(ctx context.Context, req *CreateFlowRequest) (*model.Flow, error)
	UpdateFlow(ctx context.Context, req *UpdateFlowRequest) error
	DeleteFlow(ctx context.Context, fid model.UUID) error
}

type service struct {
	g  *gorm.DB
	a  ABAC
	as audit.Service
	eb event.EventBus
	sm auth.SessionManager
}

// DeleteFlow implements Service.
func (s *service) DeleteFlow(ctx context.Context, fid model.UUID) error {
	if err := s.a.CanDeleteFlow(ctx, fid); err != nil {
		return err
	}

	if err := s.g.WithContext(ctx).Delete(&model.Flow{}, fid).Error; err != nil {
		return fmt.Errorf("delete flow by id: %w", err)
	}

	return s.as.NewAuditLog().UseContext(ctx).Action("删除流程").Payload(
		map[string]interface{}{
			"flow_id": fid,
		},
	).Commit(ctx)
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

func (s *service) CreateFlow(ctx context.Context, req *CreateFlowRequest) (*model.Flow, error) {
	session := s.sm.GetSessionFromContext(ctx)
	m := &model.Flow{
		Name:        req.Name,
		Nodes:       req.Nodes,
		Edges:       req.Edges,
		Description: req.Description,
		OwnerID:     &session.UserID,
	}
	if err := s.g.WithContext(ctx).Create(m).Error; err != nil {
		return nil, fmt.Errorf("create flow: %w", err)
	}

	s.eb.Publish(&CreateFlowEvent{
		FlowID:    m.ID,
		FlowName:  m.Name,
		OwnerID:   *m.OwnerID,
		Timestamp: time.Now().UnixMilli(),
	})

	if err := s.as.NewAuditLog().UseContext(ctx).Action("创建流程").Payload(req).Commit(ctx); err != nil {
		return nil, err
	}

	return m, nil
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

func (s *service) ListFlowByOwner(ctx context.Context, p *pagination.Pagination) ([]*model.Flow, int, error) {
	var (
		flows []*model.Flow
		count int64
	)

	u := s.sm.GetSessionFromContext(ctx)
	if err := s.g.WithContext(ctx).Model(&model.Flow{}).
		Scopes(orm.Pagination(p)).
		Select("id", "name", "owner_id", "created_at").
		Where("owner_id = ?", u.UserID).
		Find(&flows).
		Count(&count).
		Error; err != nil {
		return nil, 0, fmt.Errorf("query flow list: %w", err)
	}
	return flows, int(count), nil
}

func NewService(g *gorm.DB, sm auth.SessionManager, p ABAC, eb event.EventBus, as audit.Service) Service {
	return &service{g: g, sm: sm, a: p, as: as, eb: eb}
}

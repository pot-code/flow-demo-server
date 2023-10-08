package flow

import (
	"context"
	"errors"
	"fmt"
	"gobit-demo/infra/event"
	"gobit-demo/infra/orm"
	"gobit-demo/infra/pagination"
	"gobit-demo/model"
	"gobit-demo/services/audit"
	"gobit-demo/services/auth"
	"time"

	"gorm.io/gorm"
)

var (
	ErrDuplicatedFlow = errors.New("流程已存在")
)

type Service interface {
	GetFlowByID(ctx context.Context, id model.ID) (*model.Flow, error)
	ListFlowByOwner(ctx context.Context, p *pagination.Pagination) ([]*model.Flow, int, error)
	CreateFlow(ctx context.Context, req *CreateFlowDto) (*model.Flow, error)
	UpdateFlow(ctx context.Context, req *UpdateFlowDto) error
	DeleteFlow(ctx context.Context, id model.ID) error
}

type service struct {
	g  *gorm.DB
	a  ABAC
	as audit.Service
	eb event.EventBus
	sm auth.SessionManager
}

// DeleteFlow implements Service.
func (s *service) DeleteFlow(ctx context.Context, id model.ID) error {
	if err := s.a.CanDelete(ctx, id); err != nil {
		return err
	}

	if err := s.g.WithContext(ctx).Delete(&model.Flow{}, id).Error; err != nil {
		return fmt.Errorf("delete flow by id: %w", err)
	}

	return s.as.NewAuditLog().UseContext(ctx).Action("删除流程").Payload(
		map[string]interface{}{
			"flow_id": id,
		},
	).Commit(ctx)
}

func (s *service) GetFlowByID(ctx context.Context, id model.ID) (*model.Flow, error) {
	if err := s.a.CanView(ctx, id); err != nil {
		return nil, err
	}

	m := new(model.Flow)
	if err := s.g.WithContext(ctx).Model(&model.Flow{}).Where("id = ?", id).Take(m).Error; err != nil {
		return nil, fmt.Errorf("get flow by id: %w", err)
	}
	return m, nil
}

func (s *service) CreateFlow(ctx context.Context, req *CreateFlowDto) (*model.Flow, error) {
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

func (s *service) UpdateFlow(ctx context.Context, req *UpdateFlowDto) error {
	if err := s.a.CanUpdate(ctx, req.ID); err != nil {
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

func NewService(
	g *gorm.DB,
	sm auth.SessionManager,
	eb event.EventBus,
	as audit.Service,
) *service {
	return &service{g: g, sm: sm, a: NewPermission(g, sm), as: as, eb: eb}
}

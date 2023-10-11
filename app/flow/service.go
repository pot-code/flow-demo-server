package flow

import (
	"context"
	"errors"
	"fmt"
	"gobit-demo/infra/event"
	"gobit-demo/infra/orm"
	"gobit-demo/infra/pagination"
	"gobit-demo/model"
	"gobit-demo/model/pk"
	"gobit-demo/services/audit"
	"gobit-demo/services/auth/rbac"
	"gobit-demo/services/auth/session"
	"gobit-demo/services/notification"
	"time"

	"gorm.io/gorm"
)

var (
	ErrDuplicatedFlow = errors.New("流程已存在")
)

type Service interface {
	GetFlowByID(ctx context.Context, id pk.ID) (*model.Flow, error)
	ListFlowByOwner(ctx context.Context, p *pagination.Pagination) ([]*model.Flow, int64, error)
	CreateFlow(ctx context.Context, req *CreateFlowDto) (*model.Flow, error)
	UpdateFlow(ctx context.Context, req *UpdateFlowDto) error
	DeleteFlow(ctx context.Context, id pk.ID) error
}

type service struct {
	g  *gorm.DB
	a  *ABAC
	r  rbac.RBAC
	as audit.Service
	ns notification.Service
	eb event.EventBus
}

// DeleteFlow implements Service.
func (s *service) DeleteFlow(ctx context.Context, id pk.ID) error {
	if err := s.r.CheckPermission(ctx, "flow:delete"); err != nil {
		return err
	}
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

func (s *service) GetFlowByID(ctx context.Context, id pk.ID) (*model.Flow, error) {
	if err := s.r.CheckPermission(ctx, "flow:view"); err != nil {
		return nil, err
	}

	m := new(model.Flow)
	if err := s.g.WithContext(ctx).Model(&model.Flow{}).Where("id = ?", id).Take(m).Error; err != nil {
		return nil, fmt.Errorf("get flow by id: %w", err)
	}
	return m, nil
}

func (s *service) CreateFlow(ctx context.Context, req *CreateFlowDto) (*model.Flow, error) {
	if err := s.r.CheckPermission(ctx, "flow:create"); err != nil {
		return nil, err
	}

	session := session.GetSessionFromContext(ctx)
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

	if err := s.ns.SendNotification(ctx, *m.OwnerID, fmt.Sprintf("您创建了一个新的流程：%s", m.Name)); err != nil {
		return nil, fmt.Errorf("send notification: %w", err)
	}

	if err := s.as.NewAuditLog().UseContext(ctx).Action("创建流程").Payload(req).Commit(ctx); err != nil {
		return nil, err
	}

	return m, nil
}

func (s *service) UpdateFlow(ctx context.Context, req *UpdateFlowDto) error {
	if err := s.r.CheckPermission(ctx, "flow:update"); err != nil {
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

func (s *service) ListFlowByOwner(ctx context.Context, p *pagination.Pagination) ([]*model.Flow, int64, error) {
	if err := s.r.CheckPermission(ctx, "flow:list"); err != nil {
		return nil, -1, err
	}

	var (
		flows []*model.Flow
		count int64
	)
	u := session.GetSessionFromContext(ctx)
	if err := s.g.WithContext(ctx).Model(&model.Flow{}).
		Scopes(orm.Pagination(p)).
		Select("id", "name", "owner_id", "created_at").
		Where("owner_id = ?", u.UserID).
		Find(&flows).
		Count(&count).
		Error; err != nil {
		return nil, 0, fmt.Errorf("query flow list: %w", err)
	}
	return flows, count, nil
}

func NewService(
	g *gorm.DB,
	r rbac.RBAC,
	eb event.EventBus,
	as audit.Service,
	ns notification.Service,
) *service {
	return &service{g: g, r: r, a: NewABAC(g), as: as, eb: eb, ns: ns}
}

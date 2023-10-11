package notification

import (
	"context"
	"fmt"
	"gobit-demo/infra/orm"
	"gobit-demo/infra/pagination"
	"gobit-demo/model"
	"gobit-demo/model/pk"
	"gobit-demo/services/auth/session"

	"gorm.io/gorm"
)

type Service interface {
	SendNotification(ctx context.Context, to pk.ID, content string) error
	ListNotifications(ctx context.Context, uid pk.ID, p *pagination.Pagination) ([]*model.Notification, int64, error)
}

type service struct {
	g  *gorm.DB
	sm session.SessionManager
}

func (s *service) ListNotifications(ctx context.Context, uid pk.ID, p *pagination.Pagination) ([]*model.Notification, int64, error) {
	var (
		notifications []*model.Notification
		count         int64
	)
	if err := s.g.WithContext(ctx).
		Scopes(orm.Pagination(p)).
		Where("owner_id = ?", uid).
		Find(&notifications).
		Count(&count).Error; err != nil {
		return nil, -1, fmt.Errorf("query notification list: %w", err)
	}
	return notifications, count, nil
}

func (s *service) SendNotification(ctx context.Context, to pk.ID, content string) error {
	if err := s.g.WithContext(ctx).Create(&model.Notification{
		OwnerID: &to,
		Content: content,
	}).Error; err != nil {
		return fmt.Errorf("create notification: %w", err)
	}
	return nil
}

func NewService(g *gorm.DB, sm session.SessionManager) Service {
	return &service{g: g, sm: sm}
}

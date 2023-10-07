package notification

import (
	"context"
	"fmt"
	"gobit-demo/model"

	"gorm.io/gorm"
)

type Service interface {
	SendNotification(ctx context.Context, to model.ID, content string) error
}

type service struct {
	g *gorm.DB
}

// SendNotification implements Service.
func (s *service) SendNotification(ctx context.Context, to model.ID, content string) error {
	if err := s.g.WithContext(ctx).Create(&model.Notification{
		OwnerID: &to,
		Content: content,
	}).Error; err != nil {
		return fmt.Errorf("create notification: %w", err)
	}
	return nil
}

func NewService(g *gorm.DB) Service {
	return &service{g: g}
}

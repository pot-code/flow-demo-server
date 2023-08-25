package user

import (
	"context"
	"fmt"
	"gobit-demo/internal/pagination"
	"gobit-demo/model"
	"gobit-demo/pkg/orm"

	"gorm.io/gorm"
)

type Service interface {
	ListUser(ctx context.Context, p *pagination.Pagination) ([]*model.User, int, error)
}

type service struct {
	g *gorm.DB
}

func NewService(g *gorm.DB) *service {
	return &service{g: g}
}

func (s *service) ListUser(ctx context.Context, p *pagination.Pagination) ([]*model.User, int, error) {
	var (
		users []*model.User
		count int64
	)

	if err := s.g.WithContext(ctx).Scopes(orm.Pagination(p)).
		Select("id", "name", "username", "mobile", "disabled").
		Find(&users).
		Count(&count).
		Error; err != nil {
		return nil, 0, fmt.Errorf("query user list: %w", err)
	}

	return users, int(count), nil
}

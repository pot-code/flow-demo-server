package user

import (
	"context"
	"fmt"
	"gobit-demo/internal/pagination"
	"gobit-demo/internal/util"
	"gobit-demo/model"

	"gorm.io/gorm"
)

type Service interface {
	ListUser(ctx context.Context, p *pagination.Pagination) ([]*ListUserResponse, int, error)
}

type service struct {
	g *gorm.DB
}

func NewService(g *gorm.DB) *service {
	return &service{g: g}
}

func (s *service) ListUser(ctx context.Context, p *pagination.Pagination) ([]*ListUserResponse, int, error) {
	var (
		users []*ListUserResponse
		count int64
	)

	if err := util.GormPaginator(s.g.WithContext(ctx).Model(&model.User{}), p).
		Find(&users).
		Count(&count).
		Error; err != nil {
		return nil, 0, fmt.Errorf("query user list: %w", err)
	}

	return users, int(count), nil
}

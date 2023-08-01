package user

import (
	"context"
	"fmt"
	"gobit-demo/internal/pagination"
	"gobit-demo/model"

	"gorm.io/gorm"
)

type UserService struct {
	g *gorm.DB
}

func NewUserService(g *gorm.DB) *UserService {
	return &UserService{g: g}
}

func (s *UserService) ListUser(ctx context.Context, p *pagination.Pagination) ([]*listUserDto, uint, error) {
	var (
		users []*listUserDto
		count int64
	)

	if err := s.g.WithContext(ctx).Model(&model.User{}).
		Limit(p.PageSize).
		Offset((p.Page - 1) * p.PageSize).
		Find(&users).
		Count(&count).
		Error; err != nil {
		return nil, 0, fmt.Errorf("query user list: %w", err)
	}

	return users, uint(count), nil
}

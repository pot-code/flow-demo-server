package user

import (
	"context"
	"fmt"
	"gobit-demo/internal/pagination"
	"gobit-demo/model"

	"github.com/rs/zerolog/log"
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
		users []*model.User
		count int64
	)

	if err := s.g.WithContext(ctx).
		Limit(p.PageSize).
		Offset((p.Page - 1) * p.PageSize).
		Find(&users).Count(&count).Error; err != nil {
		return nil, 0, fmt.Errorf("query user list: %w", err)
	}

	r := make([]*listUserDto, len(users))
	for i, user := range users {
		log.Info().Uint("userId", user.ID).Msg("user")
		r[i] = new(listUserDto).fromUser(user)
	}
	return r, uint(count), nil
}

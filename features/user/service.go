package user

import (
	"context"
	"fmt"
	"gobit-demo/ent"
	"gobit-demo/internal/pagination"
)

type Service struct {
	e *ent.Client
}

func NewService(client *ent.Client) *Service {
	return &Service{e: client}
}

func (s *Service) ListUser(ctx context.Context, p *pagination.Pagination) ([]*ListUserDto, uint, error) {
	users, count, err := pagination.EntPaginator(ctx, s.e.User.Query(), p, []*ent.User{})
	if err != nil {
		return nil, 0, fmt.Errorf("query user list: %w", err)
	}

	r := make([]*ListUserDto, len(users))
	for i, user := range users {
		r[i] = new(ListUserDto).FromUser(user)
	}
	return r, uint(count), nil
}

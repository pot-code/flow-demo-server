package user

import (
	"context"
	"fmt"
	"gobit-demo/ent"
	"gobit-demo/internal/pagination"
)

type UserService struct {
	e *ent.Client
}

func NewUserService(client *ent.Client) *UserService {
	return &UserService{e: client}
}

func (s *UserService) ListUser(ctx context.Context, p *pagination.Pagination) ([]*listUserDto, uint, error) {
	users, count, err := pagination.EntPaginator(ctx, s.e.User.Query(), p, []*ent.User{})
	if err != nil {
		return nil, 0, fmt.Errorf("query user list: %w", err)
	}

	r := make([]*listUserDto, len(users))
	for i, user := range users {
		r[i] = new(listUserDto).fromUser(user)
	}
	return r, uint(count), nil
}

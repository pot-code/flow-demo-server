package user

import (
	"context"
	"gobit-demo/internal/api"
	"gobit-demo/internal/pagination"
	"net/http"
)

type UserService interface {
	ListUser(ctx context.Context, p *pagination.Pagination) ([]*ListUserDto, uint, error)
}

type controller struct {
	s UserService
}

func newController(s UserService) *controller {
	return &controller{s: s}
}

func (c *controller) list(r *http.Request, w http.ResponseWriter) error {
	p, err := api.ParsePaginationFromRequest(r)
	if err != nil {
		return err
	}

	users, count, err := c.s.ListUser(r.Context(), p)
	if err != nil {
		return err
	}

	return api.JsonPaginationResult(w, p, count, users)
}

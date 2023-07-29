package user

import (
	"context"
	"gobit-demo/internal/api"
	"gobit-demo/internal/pagination"

	"github.com/labstack/echo/v4"
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

func (c *controller) list(e echo.Context) error {
	p, err := api.ParsePaginationFromRequest(e)
	if err != nil {
		return err
	}

	users, count, err := c.s.ListUser(e.Request().Context(), p)
	if err != nil {
		return err
	}

	return api.JsonPaginationResult(e, p, count, users)
}

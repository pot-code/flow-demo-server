package user

import (
	"gobit-demo/internal/api"

	"github.com/labstack/echo/v4"
)

type controller struct {
	s Service
}

func newController(s Service) *controller {
	return &controller{s: s}
}

func (c *controller) list(e echo.Context) error {
	p, err := api.GetPaginationFromRequest(e)
	if err != nil {
		return err
	}

	users, count, err := c.s.ListUser(e.Request().Context(), p)
	if err != nil {
		return err
	}

	return api.JsonPaginationData(e, p, count, users)
}

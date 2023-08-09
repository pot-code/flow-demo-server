package user

import (
	"gobit-demo/internal/api"

	"github.com/labstack/echo/v4"
)

type route struct {
	s Service
}

func NewRoute(s Service) *route {
	return &route{s: s}
}

func (c *route) Append(g *echo.Group) {
	g.GET("", c.list)
}

func (c *route) list(e echo.Context) error {
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

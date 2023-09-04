package user

import (
	"gobit-demo/auth"
	"gobit-demo/internal/api"

	"github.com/labstack/echo/v4"
)

type route struct {
	s Service
	r auth.RBAC
}

func (c *route) Append(g *echo.Group) {
	g.GET("", c.list)
}

func (c *route) list(e echo.Context) error {
	if err := c.r.CheckPermission(e.Request().Context(), "user:list"); err != nil {
		return err
	}

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

func NewRoute(s Service, r auth.RBAC) api.Route {
	return &route{s: s, r: r}
}

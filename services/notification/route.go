package notification

import (
	"gobit-demo/infra/api"
	"gobit-demo/services/auth/session"

	"github.com/labstack/echo/v4"
)

type route struct {
	s Service
}

func (c *route) AppendRoutes(g *echo.Group) {
	g.GET("", c.list)
}

func (c *route) list(e echo.Context) error {
	p, err := api.PaginationFromRequest(e)
	if err != nil {
		return err
	}

	s := session.GetSessionFromContext(e.Request().Context())
	data, count, err := c.s.ListNotifications(e.Request().Context(), s.UserID, p)
	if err != nil {
		return err
	}
	return api.JsonPaginationData(e, p, count, data)
}

func NewRoute(s Service) *route {
	return &route{s}
}

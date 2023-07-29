package user

import (
	"gobit-demo/ent"

	"github.com/labstack/echo/v4"
)

func RegisterRoute(g *echo.Group, e *ent.Client) {
	c := newController(NewService(e))
	g.GET("/", c.list)
}

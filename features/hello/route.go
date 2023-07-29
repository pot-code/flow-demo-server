package hello

import (
	"github.com/labstack/echo/v4"
)

func RegisterRoute(g *echo.Group) {
	g.GET("/", hello)
	g.POST("/", post)
}

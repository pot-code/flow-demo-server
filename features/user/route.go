package user

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func RegisterRoute(g *echo.Group, gc *gorm.DB) {
	c := newController(NewService(gc))
	g.GET("", c.list)
}

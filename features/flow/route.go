package flow

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func RegisterRoute(g *echo.Group, gc *gorm.DB) {
	c := newController(NewFlowService(gc))
	g.POST("", c.createFlow)
	g.GET("", c.listFlow)
}

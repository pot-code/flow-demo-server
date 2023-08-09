package flow

import (
	"gobit-demo/features/perm"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func RegisterRoute(g *echo.Group, db *gorm.DB, ps perm.Service) {
	c := newController(NewService(db), ps)
	g.POST("", c.createFlow)
	g.GET("", c.listFlow)
	g.GET("/node", c.listFlowNode)
	g.POST("/node", c.createFlowNode)
}

package flow

import (
	"gobit-demo/features/perm"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func RegisterRoute(g *echo.Group, db *gorm.DB) {
	c := newController(NewService(db), perm.NewService(db))
	g.POST("", c.createFlow)
	g.GET("", c.listFlow)
	g.GET("/node", c.listFlowNode)
	g.POST("/node", c.createFlowNode)
}

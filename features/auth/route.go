package auth

import (
	"gobit-demo/internal/event"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func RegisterRoute(g *echo.Group, gc *gorm.DB, eb event.EventBus, ts tokenService, exp time.Duration) {
	c := newController(NewAuthService(gc, eb), ts)
	g.POST("/login", c.login)
	g.PUT("/logout", c.logout)
	g.POST("/register", c.register)
}

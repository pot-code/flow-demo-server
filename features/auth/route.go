package auth

import (
	"gobit-demo/internal/event"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func RegisterRoute(g *echo.Group, gc *gorm.DB, eb event.EventBus, ts TokenService) {
	c := newController(NewService(gc, eb, NewPasswordHash()), ts)
	g.POST("/login", c.login)
	g.PUT("/logout", c.logout)
	g.POST("/register", c.register)
}

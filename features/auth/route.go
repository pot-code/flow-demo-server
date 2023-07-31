package auth

import (
	"gobit-demo/internal/token"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func RegisterRoute(g *echo.Group, gc *gorm.DB, jwt *token.JwtIssuer, exp time.Duration) {
	c := newController(NewAuthService(gc), NewJwtService(jwt, exp))
	g.POST("/login", c.login)
	g.POST("/register", c.register)
}

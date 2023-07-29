package auth

import (
	"gobit-demo/ent"
	"gobit-demo/internal/token"
	"time"

	"github.com/labstack/echo/v4"
)

func RegisterRoute(g *echo.Group, e *ent.Client, jwt *token.JwtIssuer, exp time.Duration) {
	c := newController(NewService(e, jwt, exp))
	g.POST("/login", c.login)
	g.POST("/register", c.register)
}

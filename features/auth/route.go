package auth

import (
	"gobit-demo/internal/token"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func RegisterRoute(g *echo.Group, gc *gorm.DB, jwt *token.JwtIssuer, rc *redis.Client, exp time.Duration) {
	jb := newRedisBlacklist(rc, exp)
	c := newController(NewAuthService(gc), NewJwtService(jwt, jb, exp))
	g.POST("/login", c.login)
	g.PUT("/logout", c.logout)
	g.POST("/register", c.register)
}

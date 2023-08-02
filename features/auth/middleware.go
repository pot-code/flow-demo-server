package auth

import (
	"gobit-demo/internal/api"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func getJwtTokenFromRequest(c echo.Context) string {
	value := c.Request().Header.Get("Authorization")
	if strings.HasPrefix(value, "Bearer ") {
		return value[7:]
	}
	return ""
}

type jwtService interface {
	Verify(token string) (jwt.Claims, error)
	ClaimToUser(claims jwt.Claims) *LoginUser
}

func AuthMiddleware(ts jwtService) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := getJwtTokenFromRequest(c)
			if token == "" {
				return api.JsonUnauthenticated(c, "未登录")
			}

			claim, err := ts.Verify(token)
			if err != nil {
				return api.JsonUnauthorized(c, "token 无效")
			}

			u := ts.ClaimToUser(claim)
			c.SetRequest(c.Request().WithContext(setLoginUserContextValue(c.Request().Context(), u)))
			return next(c)
		}
	}
}

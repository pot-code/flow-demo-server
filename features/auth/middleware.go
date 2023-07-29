package auth

import (
	"gobit-demo/internal/api"
	"gobit-demo/internal/token"
	"strings"

	"github.com/labstack/echo/v4"
)

func getJwtTokenFromRequest(c echo.Context) string {
	value := c.Request().Header.Get("Authorization")
	if strings.HasPrefix(value, "Bearer ") {
		return value[7:]
	}
	return ""
}

func AuthMiddleware(jwt *token.JwtIssuer) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := getJwtTokenFromRequest(c)
			if token == "" {
				return api.JsonUnauthenticated(c, "没有登录")
			}

			claim, err := jwt.Verify(token)
			if err != nil {
				return api.JsonUnauthorized(c, "token 无效")
			}

			u := new(LoginUser).fromClaim(claim)
			c.SetRequest(c.Request().WithContext(setUserContextValue(c.Request().Context(), u)))
			return next(c)
		}
	}
}

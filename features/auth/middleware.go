package auth

import (
	"gobit-demo/internal/api"
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

func AuthMiddleware(ts TokenService) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := getJwtTokenFromRequest(c)
			if token == "" {
				return api.JsonUnauthenticated(c, "未登录")
			}

			u, err := ts.Verify(token)
			if err != nil {
				return api.JsonUnauthorized(c, "token 无效")
			}

			ok, err := ts.IsInBlacklist(c.Request().Context(), token)
			if err != nil {
				return err
			}
			if ok {
				return api.JsonUnauthorized(c, "token 无效")
			}

			c.SetRequest(c.Request().WithContext(u.WithContext(c.Request().Context())))
			return next(c)
		}
	}
}

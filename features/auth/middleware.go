package auth

import (
	"errors"
	"gobit-demo/internal/api"
	"strings"

	"github.com/labstack/echo/v4"
)

func getTokenFromRequest(c echo.Context) string {
	value := c.Request().Header.Get("Authorization")
	if strings.HasPrefix(value, "Bearer ") {
		return value[7:]
	}
	return ""
}

func AuthMiddleware(ts TokenService, sm SessionManager) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := getTokenFromRequest(c)
			if token == "" {
				return api.JsonUnauthenticated(c, "未登录")
			}

			u, err := ts.Verify(token)
			if err != nil {
				return api.JsonUnauthorized(c, "token 无效")
			}

			s, err := sm.GetSession(c.Request().Context(), u.SessionID)
			if errors.Is(err, ErrSessionNotFound) {
				return api.JsonUnauthorized(c, "token 无效")
			}
			if err != nil {
				return api.JsonServerError(c, err.Error())
			}

			c.SetRequest(c.Request().WithContext(sm.SetSession(c.Request().Context(), s)))
			return next(c)
		}
	}
}

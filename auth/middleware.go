package auth

import (
	"errors"
	"gobit-demo/internal/api"
	"time"

	"github.com/labstack/echo/v4"
)

func AuthMiddleware(ts TokenService, sm SessionManager, threshold time.Duration) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token, _ := ts.FromHttpRequest(c.Request())
			if token == "" {
				return api.JsonUnauthorized(c, "未登录")
			}
			u, err := ts.Verify(token)
			if err != nil {
				return api.JsonUnauthorized(c, "token 无效")
			}

			s, err := sm.GetSessionBySessionID(c.Request().Context(), u.SessionID)
			if errors.Is(err, ErrSessionNotFound) {
				return api.JsonUnauthorized(c, "token 无效")
			}
			if err != nil {
				return api.JsonServerError(c, err.Error())
			}

			if err := sm.RefreshSession(c.Request().Context(), s, threshold); err != nil {
				return api.JsonServerError(c, err.Error())
			}

			c.SetRequest(c.Request().WithContext(sm.SetSession(c.Request().Context(), s)))
			return next(c)
		}
	}
}

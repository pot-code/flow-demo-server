package middlewares

import (
	"errors"
	"gobit-demo/infra/api"
	"gobit-demo/services/auth"
	"gobit-demo/services/auth/session"
	"time"

	"github.com/labstack/echo/v4"
)

func AuthMiddleware(ts auth.TokenService, sm session.SessionManager, threshold time.Duration) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token, _ := ts.GetTokenFromRequest(c.Request())
			if token == "" {
				return api.JsonUnauthorized(c, "未登录")
			}
			u, err := ts.Verify(token)
			if err != nil {
				return api.JsonUnauthorized(c, "token 无效")
			}

			s, err := sm.GetSessionBySessionID(c.Request().Context(), u.SessionID)
			if errors.Is(err, session.ErrSessionNotFound) {
				return api.JsonUnauthorized(c, "token 无效")
			}
			if err != nil {
				return api.JsonServerError(c, err.Error())
			}

			if err := sm.RefreshSession(c.Request().Context(), s, threshold); err != nil {
				return api.JsonServerError(c, err.Error())
			}

			c.SetRequest(c.Request().WithContext(session.WithSessionContext(c.Request().Context(), s)))
			return next(c)
		}
	}
}

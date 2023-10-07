package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
)

var LoggingMiddleware = middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
	LogURI:    true,
	LogStatus: true,
	LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
		log.Debug().
			Str("method", c.Request().Method).
			Str("uri", v.URI).
			Int("status", v.Status).
			Dur("latency", v.Latency).
			Msg("request")

		return v.Error
	},
})

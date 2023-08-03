package main

import (
	"fmt"
	"gobit-demo/features/auth"
	"gobit-demo/features/flow"
	"gobit-demo/features/user"
	"gobit-demo/internal/api"
	"gobit-demo/internal/cache"
	"gobit-demo/internal/config"
	"gobit-demo/internal/db"
	"gobit-demo/internal/logging"
	"gobit-demo/internal/token"
	"gobit-demo/internal/validate"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func main() {
	validate.Init()
	cfg := config.LoadConfig()
	logging.Init(cfg)

	log.Debug().Any("config", cfg).Msg("config")

	rc := cache.NewRedisCache(cfg.Cache.DSN)
	dc := db.NewDB(cfg.Database.DSN)
	gc, err := db.NewGormClient(dc)
	if err != nil {
		log.Fatal().Err(err).Msg("error creating gorm client")
	}
	jwt := token.NewJwtIssuer(cfg.Token.Secret)

	e := echo.New()
	e.HTTPErrorHandler = api.ErrorHandler
	e.Use(api.LoggingMiddleware)

	api.GroupRoute(e, "/auth", func(g *echo.Group) {
		auth.RegisterRoute(g, gc, jwt, rc, cfg.Token.Exp)
	})
	api.GroupRoute(e, "/flow", func(g *echo.Group) {
		g.Use(auth.AuthMiddleware(jwt))
		flow.RegisterRoute(g, gc)
	})
	api.GroupRoute(e, "/user", func(g *echo.Group) {
		g.Use(auth.AuthMiddleware(jwt))
		user.RegisterRoute(g, gc)
	})

	if err := e.Start(fmt.Sprintf(":%d", cfg.Port)); err != http.ErrServerClosed {
		log.Err(err).Msg("error starting server")
	}
}

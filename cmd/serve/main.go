package main

import (
	"fmt"
	"gobit-demo/config"
	"gobit-demo/features/auth"
	"gobit-demo/features/flow"
	"gobit-demo/features/user"
	"gobit-demo/internal/api"
	"gobit-demo/internal/cache"
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
	logging.Init(cfg.Logging.Level)

	log.Debug().Any("config", cfg).Msg("config")

	rc := cache.NewRedisCache(cfg.Cache.DSN)
	dc := db.NewDB(cfg.Database.DSN)
	gc, err := db.NewGormClient(dc, log.Logger)
	if err != nil {
		log.Fatal().Err(err).Msg("error creating gorm client")
	}
	js := auth.NewJwtService(
		token.NewJwtIssuer(cfg.Token.Secret),
		auth.NewRedisBlacklist(rc, cfg.Token.Exp),
		cfg.Token.Exp,
	)

	e := echo.New()
	e.HTTPErrorHandler = api.ErrorHandler
	e.Use(api.LoggingMiddleware)

	api.GroupRoute(e, "/auth", func(g *echo.Group) {
		auth.RegisterRoute(g, gc, js, cfg.Token.Exp)
	})
	api.GroupRoute(e, "/flow", func(g *echo.Group) {
		g.Use(auth.AuthMiddleware(js))
		flow.RegisterRoute(g, gc)
	})
	api.GroupRoute(e, "/user", func(g *echo.Group) {
		g.Use(auth.AuthMiddleware(js))
		user.RegisterRoute(g, gc)
	})

	if err := e.Start(fmt.Sprintf(":%d", cfg.Port)); err != http.ErrServerClosed {
		log.Err(err).Msg("error starting server")
	}
}

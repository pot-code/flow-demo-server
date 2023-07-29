package main

import (
	"fmt"
	"gobit-demo/features/auth"
	"gobit-demo/features/hello"
	"gobit-demo/features/user"
	"gobit-demo/internal/api"
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

	conn := db.NewDB(cfg.Database.DSN)
	ent := db.NewEntClient(conn)
	jwt := token.NewJwtIssuer(cfg.Token.Secret)

	e := echo.New()
	e.HTTPErrorHandler = api.ErrorHandler
	e.Use(api.LoggingMiddleware)

	api.GroupRoute(e, "/auth", func(g *echo.Group) {
		auth.RegisterRoute(g, ent, jwt, cfg.Token.Exp)
	})
	api.GroupRoute(e, "/hello", func(g *echo.Group) {
		g.Use(auth.AuthMiddleware(jwt))
		hello.RegisterRoute(g)
	})
	api.GroupRoute(e, "/user", func(g *echo.Group) {
		g.Use(auth.AuthMiddleware(jwt))
		user.RegisterRoute(g, ent)
	})

	if err := e.Start(fmt.Sprintf(":%d", cfg.Port)); err != http.ErrServerClosed {
		log.Err(err).Msg("error starting server")
	}
}

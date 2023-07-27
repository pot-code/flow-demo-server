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

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

func main() {
	validate.Init()
	cfg := config.LoadConfig()
	logging.Init(cfg)

	d := db.NewDB(cfg.Database.DSN)
	e := db.NewEntClient(d)
	j := token.NewJwtIssuer(cfg.Token.Secret)

	mux := chi.NewRouter()
	mux.Use(api.LoggingMiddleware)

	mux.Route("/auth", func(r chi.Router) {
		auth.RegisterRoute(r, e, j, cfg.Token.Exp)
	})
	mux.Route("/hello", func(r chi.Router) {
		r.Use(auth.AuthMiddleware(j))
		hello.RegisterRoute(r)
	})
	mux.Route("/user", func(r chi.Router) {
		r.Use(auth.AuthMiddleware(j))
		user.RegisterRoute(r, e)
	})

	log.Info().Int("port", cfg.Port).Msg("starting server")
	if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), mux); err != nil {
		log.Err(err).Msg("error starting server")
	}
}

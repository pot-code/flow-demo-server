package main

import (
	"context"
	"gobit-demo/ent/migrate"
	"gobit-demo/internal/config"
	"gobit-demo/internal/db"
	"gobit-demo/internal/logging"
	"gobit-demo/internal/validate"

	"github.com/rs/zerolog/log"
)

func main() {
	validate.Init()
	cfg := config.LoadConfig()
	logging.Init(cfg)

	d := db.NewDB(cfg.Database.DSN)
	e := db.NewEntClient(d)

	if err := e.Schema.Create(context.Background(),
		migrate.WithDropColumn(true),
		migrate.WithDropIndex(true),
	); err != nil {
		log.Err(err).Msg("error migrating schema")
	}
}

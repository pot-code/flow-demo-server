package main

import (
	"gobit-demo/internal/config"
	"gobit-demo/internal/db"
	"gobit-demo/internal/logging"
	"gobit-demo/internal/validate"
	"gobit-demo/model"

	"github.com/rs/zerolog/log"
)

func main() {
	validate.Init()
	cfg := config.LoadConfig()
	logging.Init(cfg)

	d := db.NewDB(cfg.Database.DSN)
	g, err := db.NewGormClient(d)
	if err != nil {
		log.Fatal().Err(err).Msg("error creating gorm client")
	}

	if err := g.AutoMigrate(&model.User{}); err != nil {
		log.Fatal().Err(err).Msg("error migrating schema")
	}
}

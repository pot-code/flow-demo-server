package main

import (
	"gobit-demo/config"
	"gobit-demo/internal/db"
	"gobit-demo/internal/logging"
	"gobit-demo/internal/orm"
	"gobit-demo/internal/validate"
	"gobit-demo/model"

	"github.com/rs/zerolog/log"
)

func main() {
	validate.Init()
	cfg := config.LoadConfig()
	logging.Init(cfg.Logging.Level)

	d := db.NewMysqlDB(cfg.Database.String())
	gc := orm.NewGormDB(d, log.Logger)

	if err := gc.AutoMigrate(
		&model.CasbinRule{},
		&model.User{},
		&model.Role{},
		&model.Flow{},
		&model.FlowNode{},
	); err != nil {
		log.Fatal().Err(err).Msg("error migrating schema")
	}
}

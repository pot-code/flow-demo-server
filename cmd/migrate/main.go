package main

import (
	"gobit-demo/config"
	"gobit-demo/internal/db"
	"gobit-demo/internal/logging"
	"gobit-demo/internal/orm"
	"gobit-demo/model"

	"github.com/rs/zerolog/log"
)

func main() {
	cfg := config.LoadConfig()
	logging.Init(cfg.Logging.Level)

	d := db.NewMysqlDB(cfg.Database.GetDSN())
	gc := orm.NewGormDB(d, log.Logger)

	if err := gc.AutoMigrate(
		&model.AuditLog{},
		&model.User{},
		&model.Role{},
		&model.Permission{},
		&model.Flow{},
	); err != nil {
		log.Fatal().Err(err).Msg("error migrating schema")
	}
}

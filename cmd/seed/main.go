package main

import (
	"context"
	"gobit-demo/config"
	"gobit-demo/infra/db"
	"gobit-demo/infra/logging"
	"gobit-demo/infra/orm"
	"gobit-demo/model"

	"github.com/rs/zerolog/log"
)

func main() {
	cfg := config.LoadConfig()
	logging.Init(cfg.Logging.Level)

	d := db.NewMysqlDB(cfg.Database.GetDSN())
	g := orm.NewGormDB(d, log.Logger)

	permissions := []*model.Permission{
		{
			Name: "flow:list",
		},
		{
			Name: "flow:view",
		},
		{
			Name: "flow:create",
		},
		{
			Name: "flow:update",
		},
		{
			Name: "flow:delete",
		},
	}
	if err := g.WithContext(context.Background()).CreateInBatches(permissions, len(permissions)).Error; err != nil {
		log.Err(err).Msg("seed permissions")
	}

	roles := []*model.Role{
		{
			Name: "admin",
		},
		{
			Name:        "user",
			Permissions: permissions,
		},
	}
	if err := g.WithContext(context.Background()).CreateInBatches(roles, len(roles)).Error; err != nil {
		log.Err(err).Msg("seed roles")
	}
}

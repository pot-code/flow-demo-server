package main

import (
	"context"
	"gobit-demo/config"
	"gobit-demo/features/auth"
	"gobit-demo/internal/db"
	"gobit-demo/internal/logging"
	"gobit-demo/internal/orm"
	"gobit-demo/internal/uuid"
	"gobit-demo/model"

	"github.com/rs/zerolog/log"
)

func main() {
	cfg := config.LoadConfig()
	logging.Init(cfg.Logging.Level)
	uuid.InitSonyflake(cfg.NodeID)

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
			Name: "flow:list",
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

	us := auth.NewService(g)
	user := &auth.CreateUserRequest{
		Name:     cfg.Admin.Name,
		Username: cfg.Admin.Username,
		Password: cfg.Admin.Password,
	}
	m, err := us.CreateUser(context.Background(), user)
	if err != nil {
		log.Err(err).Msg("create admin user")
	}
	m.Roles = append(m.Roles, roles[0])
	if err := g.WithContext(context.Background()).Save(m).Error; err != nil {
		log.Err(err).Msg("assign admin role")
	}
}
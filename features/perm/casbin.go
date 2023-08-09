package perm

import (
	"fmt"
	"gobit-demo/model"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

func NewCasbinEnforcer(gd *gorm.DB) *casbin.Enforcer {
	gormadapter.TurnOffAutoMigrate(gd)
	a, err := gormadapter.NewAdapterByDBWithCustomTable(gd, &model.CasbinRule{})
	if err != nil {
		panic(fmt.Errorf("error creating casbin gorm adapter: %w", err))
	}

	e, err := casbin.NewEnforcer("casbin/model.conf", a)
	if err != nil {
		panic(fmt.Errorf("error creating casbin enforcer: %w", err))
	}
	e.EnableAutoSave(true)
	return e
}

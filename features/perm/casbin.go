package perm

import (
	"fmt"
	"gobit-demo/model"
	"strconv"

	"github.com/casbin/casbin/v2/log"
	"github.com/casbin/casbin/v2/rbac"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

func NewCasbinEnforcer(gd *gorm.DB) *casbin.Enforcer {
	gormadapter.TurnOffAutoMigrate(gd)
	a, err := gormadapter.NewAdapterByDBWithCustomTable(gd, &model.CasbinRule{}, "casbin_rules")
	if err != nil {
		panic(fmt.Errorf("error creating casbin gorm adapter: %w", err))
	}

	e, err := casbin.NewEnforcer("casbin/model.conf", a)
	if err != nil {
		panic(fmt.Errorf("error creating casbin enforcer: %w", err))
	}
	e.SetRoleManager(newGormRoleManager(gd))
	e.EnableAutoSave(true)
	e.LoadPolicy()
	return e
}

type gormRoleManager struct {
	g *gorm.DB
}

func newGormRoleManager(g *gorm.DB) *gormRoleManager {
	return &gormRoleManager{g: g}
}

func (r *gormRoleManager) Clear() error {
	return nil
}

func (r *gormRoleManager) AddLink(name1 string, name2 string, domain ...string) error {
	panic("implement me")
}

func (r *gormRoleManager) BuildRelationship(name1 string, name2 string, domain ...string) error {
	panic("implement me")
}

func (r *gormRoleManager) DeleteLink(name1 string, name2 string, domain ...string) error {
	panic("implement me")
}

func (r *gormRoleManager) HasLink(name1 string, name2 string, domain ...string) (bool, error) {
	roles, err := r.GetRoles(name1, domain...)
	if err != nil {
		return false, err
	}
	for _, role := range roles {
		if role == name2 {
			return true, nil
		}
	}
	return false, nil
}

func (r *gormRoleManager) GetRoles(name string, domain ...string) ([]string, error) {
	uid, err := strconv.Atoi(name)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}

	var roles []string
	if err := r.g.Table("user_roles").
		Select("roles.name").
		Joins("INNER JOIN roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", uid).
		Scan(&roles).Error; err != nil {
		return nil, fmt.Errorf("get user roles: %w", err)
	}
	return roles, nil
}

func (r *gormRoleManager) GetUsers(name string, domain ...string) ([]string, error) {
	rid, err := strconv.Atoi(name)
	if err != nil {
		return nil, fmt.Errorf("invalid role id: %w", err)
	}

	var users []string
	if err := r.g.Table("user_roles").
		Select("users.id").
		Joins("INNER JOIN roles ON user_roles.user_id = users.id").
		Where("user_roles.role_id = ?", rid).
		Scan(&users).Error; err != nil {
		return nil, fmt.Errorf("get user roles: %w", err)
	}
	return users, nil
}

func (r *gormRoleManager) GetDomains(name string) ([]string, error) {
	panic("implement me")
}

func (r *gormRoleManager) GetAllDomains() ([]string, error) {
	panic("implement me")
}

func (r *gormRoleManager) PrintRoles() error {
	panic("implement me")
}

func (r *gormRoleManager) SetLogger(logger log.Logger) {
	panic("implement me")
}

func (r *gormRoleManager) Match(str string, pattern string) bool {
	panic("implement me")
}

func (r *gormRoleManager) AddMatchingFunc(name string, fn rbac.MatchingFunc) {
	panic("implement me")
}

func (r *gormRoleManager) AddDomainMatchingFunc(name string, fn rbac.MatchingFunc) {
	panic("implement me")
}

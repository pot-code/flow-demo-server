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

func NewCasbinEnforcer(db *gorm.DB) *casbin.Enforcer {
	gormadapter.TurnOffAutoMigrate(db)
	a, err := gormadapter.NewAdapterByDBWithCustomTable(db, &model.CasbinRule{}, "casbin_rules")
	if err != nil {
		panic(fmt.Errorf("error creating casbin gorm adapter: %w", err))
	}

	e, err := casbin.NewEnforcer("casbin/model.conf", a)
	if err != nil {
		panic(fmt.Errorf("error creating casbin enforcer: %w", err))
	}
	e.SetRoleManager(newRoleManager(db))
	e.EnableAutoSave(true)
	e.LoadPolicy()
	return e
}

func newRoleManager(g *gorm.DB) rbac.RoleManager {
	return &roleManager{g: g}
}

type roleManager struct {
	g *gorm.DB
}

func (r *roleManager) Clear() error {
	return nil
}

func (r *roleManager) AddLink(name1 string, name2 string, domain ...string) error {
	panic("implement me")
}

func (r *roleManager) BuildRelationship(name1 string, name2 string, domain ...string) error {
	panic("implement me")
}

func (r *roleManager) DeleteLink(name1 string, name2 string, domain ...string) error {
	panic("implement me")
}

func (r *roleManager) HasLink(name1 string, name2 string, domain ...string) (bool, error) {
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

func (r *roleManager) GetRoles(name string, domain ...string) ([]string, error) {
	uid, err := r.parseUserId(name)
	if err != nil {
		return nil, err
	}
	return r.getUserRoles(uid)
}

func (r *roleManager) getUserRoles(uid interface{}) ([]string, error) {
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

func (r *roleManager) parseUserId(name string) (interface{}, error) {
	uid, err := strconv.Atoi(name)
	if err != nil {
		return 0, fmt.Errorf("invalid user id: %w", err)
	}
	return uint(uid), nil
}

func (r *roleManager) GetUsers(name string, domain ...string) ([]string, error) {
	rid, err := r.parseRoleId(name)
	if err != nil {
		return nil, err
	}

	var users []string
	if err := r.g.Table("user_roles").
		Select("users.id").
		Joins("INNER JOIN roles ON user_roles.user_id = users.id").
		Where("user_roles.role_id = ?", rid).
		Scan(&users).Error; err != nil {
		return nil, fmt.Errorf("get role users: %w", err)
	}
	return users, nil
}

func (r *roleManager) parseRoleId(name string) (interface{}, error) {
	uid, err := strconv.Atoi(name)
	if err != nil {
		return 0, fmt.Errorf("invalid role id: %w", err)
	}
	return uint(uid), nil
}

func (r *roleManager) GetDomains(name string) ([]string, error) {
	panic("implement me")
}

func (r *roleManager) GetAllDomains() ([]string, error) {
	panic("implement me")
}

func (r *roleManager) PrintRoles() error {
	panic("implement me")
}

func (r *roleManager) SetLogger(logger log.Logger) {
	panic("implement me")
}

func (r *roleManager) Match(str string, pattern string) bool {
	panic("implement me")
}

func (r *roleManager) AddMatchingFunc(name string, fn rbac.MatchingFunc) {
	panic("implement me")
}

func (r *roleManager) AddDomainMatchingFunc(name string, fn rbac.MatchingFunc) {
	panic("implement me")
}

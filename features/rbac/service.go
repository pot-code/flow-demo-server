package rbac

import (
	"context"
	"fmt"
	"gobit-demo/features/auth"

	"github.com/open-policy-agent/opa/rego"
	"gorm.io/gorm"
)

type Service interface {
	HasPermission(ctx context.Context, obj, act string) error
}

func NewService(g *gorm.DB) Service {
	query, err := rego.New(
		rego.Query("data.rbac.allow"),
		rego.Load([]string{"opa/rbac.rego"}, nil),
	).PrepareForEval(context.Background())
	if err != nil {
		panic(fmt.Errorf("create rego query: %w", err))
	}
	return &service{g: g, q: &query}
}

type service struct {
	g *gorm.DB
	q *rego.PreparedEvalQuery
}

func (s *service) HasPermission(ctx context.Context, obj string, act string) error {
	u, ok := new(auth.LoginUser).FromContext(ctx)
	if !ok {
		panic(fmt.Errorf("no login user attached in context"))
	}

	res, err := s.q.Eval(ctx, rego.EvalInput(map[string]any{
		"sub": u.ID,
		"obj": obj,
		"act": act,
	}))
	if err != nil {
		return fmt.Errorf("evaluate policy: %w", err)
	}
	if !res.Allowed() {
		return &NoPermissionError{
			UserID:   u.ID,
			Username: u.Username,
			Obj:      obj,
			Act:      act,
		}
	}
	return nil
}

type PolicyService interface {
	AddPolicy(ctx context.Context, role, obj, act string) error
	UpdatePolicy(ctx context.Context, role, obj, act string) error
	DeletePolicy(ctx context.Context, role, obj, act string) error
}

func NewPolicyService(g *gorm.DB) PolicyService {
	return &policyService{g: g}
}

type policyService struct {
	g *gorm.DB
}

// AddPolicy implements PolicyService.
func (*policyService) AddPolicy(ctx context.Context, role string, obj string, act string) error {
	panic("unimplemented")
}

// DeletePolicy implements PolicyService.
func (*policyService) DeletePolicy(ctx context.Context, role string, obj string, act string) error {
	panic("unimplemented")
}

// UpdatePolicy implements PolicyService.
func (*policyService) UpdatePolicy(ctx context.Context, role string, obj string, act string) error {
	panic("unimplemented")
}

type RoleService interface {
	AddRole(ctx context.Context, role string) error
	UpdateRole(ctx context.Context, role string) error
	DeleteRole(ctx context.Context, role string) error
}

func NewRoleService(g *gorm.DB) RoleService {
	return &roleService{g: g}
}

type roleService struct {
	g *gorm.DB
}

// AddRole implements RoleService.
func (*roleService) AddRole(ctx context.Context, role string) error {
	panic("unimplemented")
}

// DeleteRole implements RoleService.
func (*roleService) DeleteRole(ctx context.Context, role string) error {
	panic("unimplemented")
}

// UpdateRole implements RoleService.
func (*roleService) UpdateRole(ctx context.Context, role string) error {
	panic("unimplemented")
}

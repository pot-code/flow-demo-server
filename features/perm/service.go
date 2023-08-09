package perm

import (
	"context"
	"fmt"
	"gobit-demo/features/auth"
	"strconv"

	"github.com/casbin/casbin/v2"
)

type Service interface {
	HasPerm(ctx context.Context, obj, act string) error
	AddPerm(ctx context.Context, role, obj, act string) error
	DeletePerm(ctx context.Context, role, obj, act string) error
}

func NewService(e *casbin.Enforcer) Service {
	return &service{e: e}
}

type service struct {
	e *casbin.Enforcer
}

func (s *service) AddPerm(ctx context.Context, role string, obj string, act string) error {
	if _, err := s.e.AddPolicy(role, obj, act); err != nil {
		return fmt.Errorf("add permission: %w", err)
	}
	return nil
}

func (s *service) DeletePerm(ctx context.Context, role string, obj string, act string) error {
	if _, err := s.e.RemovePolicy(role, obj, act); err != nil {
		return fmt.Errorf("delete permission: %w", err)
	}
	return nil
}

func (s *service) HasPerm(ctx context.Context, obj string, act string) error {
	u, ok := auth.GetLoginUserFromContext(ctx)
	if !ok {
		return &NoPermissionError{
			Obj: obj,
			Act: act,
		}
	}

	ok, err := s.e.Enforce(strconv.Itoa(int(u.Id)), obj, act)
	if err != nil {
		return fmt.Errorf("check permission: %w", err)
	}
	if !ok {
		return &NoPermissionError{
			UserID:   u.Id,
			Username: u.Username,
			Obj:      obj,
			Act:      act,
		}
	}
	return nil
}

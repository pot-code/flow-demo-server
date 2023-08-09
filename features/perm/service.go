package perm

import (
	"context"
	"gobit-demo/features/auth"

	"github.com/casbin/casbin/v2"
	"github.com/rs/zerolog/log"
)

type PermService interface {
	HasPerm(ctx context.Context, obj, act string) bool
}

type casbinPermService struct {
	e *casbin.Enforcer
}

func NewCasbinPermService(e *casbin.Enforcer) *casbinPermService {
	return &casbinPermService{e: e}
}

func (s *casbinPermService) HasPerm(ctx context.Context, obj, act string) bool {
	sub := auth.GetLoginUserFromContext(ctx)
	ok, err := s.e.Enforce(sub.Id, obj, act)
	if err != nil {
		log.Err(err).
			Uint("sub", sub.Id).
			Str("obj", obj).
			Str("act", act).
			Msg("error checking permission")
	}
	return ok
}

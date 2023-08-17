package auth

import (
	"context"
)

type contextKey struct{}

var (
	sessionKey = contextKey{}
)

func (s *Session) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, sessionKey, s)
}

func (s *Session) FromContext(ctx context.Context) (*Session, bool) {
	v, ok := ctx.Value(sessionKey).(*Session)
	return v, ok
}

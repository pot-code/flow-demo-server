package session

import "context"

type sessionKeyType struct{}

var sessionKey = sessionKeyType{}

func GetSessionFromContext(ctx context.Context) *Session {
	v, ok := ctx.Value(sessionKey).(*Session)
	if !ok {
		panic("session not found in context")
	}
	return v
}

func WithSession(ctx context.Context, s *Session) context.Context {
	return context.WithValue(ctx, sessionKey, s)
}

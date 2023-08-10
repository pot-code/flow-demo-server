package auth

import (
	"context"
)

type contextKey struct{}

var userKey = contextKey{}

func (u *LoginUser) FromContext(ctx context.Context) (*LoginUser, bool) {
	v, ok := ctx.Value(userKey).(*LoginUser)
	return v, ok
}

func (u *LoginUser) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, userKey, u)
}

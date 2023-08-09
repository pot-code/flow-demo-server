package auth

import (
	"context"
)

type authKeyType struct{}

var userKey = authKeyType{}

func (u *LoginUser) FromContext(ctx context.Context) (*LoginUser, bool) {
	v, ok := ctx.Value(userKey).(*LoginUser)
	return v, ok
}

func (u *LoginUser) ToContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, userKey, u)
}

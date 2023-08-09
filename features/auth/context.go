package auth

import (
	"context"
)

type authKeyType struct{}

var userKey = authKeyType{}

func GetLoginUserFromContext(ctx context.Context) (*LoginUser, bool) {
	v, ok := ctx.Value(userKey).(*LoginUser)
	return v, ok
}

func setLoginUserContextValue(ctx context.Context, u *LoginUser) context.Context {
	return context.WithValue(ctx, userKey, u)
}

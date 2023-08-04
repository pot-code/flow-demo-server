package auth

import (
	"context"
)

type authKeyType struct{}

var userKey = authKeyType{}

func GetLoginUserFromContext(ctx context.Context) *LoginUser {
	return ctx.Value(userKey).(*LoginUser)
}

func setLoginUserContextValue(ctx context.Context, u *LoginUser) context.Context {
	return context.WithValue(ctx, userKey, u)
}

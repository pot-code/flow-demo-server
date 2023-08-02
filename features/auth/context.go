package auth

import (
	"context"
)

type authKey string

func GetLoginUserFromContext(ctx context.Context) *LoginUser {
	return ctx.Value(authKey("user")).(*LoginUser)
}

func setLoginUserContextValue(ctx context.Context, u *LoginUser) context.Context {
	return context.WithValue(ctx, authKey("user"), u)
}

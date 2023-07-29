package auth

import (
	"context"
)

type userKey string

func UserFromContext(ctx context.Context) *LoginUser {
	return ctx.Value(userKey("user")).(*LoginUser)
}

func setUserContextValue(ctx context.Context, u *LoginUser) context.Context {
	return context.WithValue(ctx, userKey("user"), u)
}

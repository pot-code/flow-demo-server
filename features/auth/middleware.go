package auth

import (
	"gobit-demo/internal/api"
	"gobit-demo/internal/token"
	"net/http"
	"strings"
)

func getJwtTokenFromRequest(r *http.Request) string {
	value := r.Header.Get("Authorization")
	if strings.HasPrefix(value, "Bearer ") {
		return value[7:]
	}
	return ""
}

func AuthMiddleware(jwt *token.JwtIssuer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := getJwtTokenFromRequest(r)
			if token == "" {
				api.JsonUnauthenticated(w, "没有登录")
				return
			}

			c, err := jwt.Verify(token)
			if err != nil {
				api.JsonUnauthorized(w, "token 无效")
				return
			}

			u := new(LoginUser).fromClaim(c)
			next.ServeHTTP(
				w,
				r.WithContext(u.setContextValue(r.Context())),
			)
		})
	}
}

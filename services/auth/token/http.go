package token

import (
	"fmt"
	"net/http"
)

type HttpTokenHelper interface {
	GetTokenFromRequest(r *http.Request) (string, error)
	SetTokenInResponse(w http.ResponseWriter, token string)
}

type cookieHelper struct {
	key string
}

func NewHttpCookieTokenHelper(key string) HttpTokenHelper {
	return &cookieHelper{
		key: key,
	}
}

func (c *cookieHelper) GetTokenFromRequest(r *http.Request) (string, error) {
	token, err := r.Cookie(c.key)
	if err != nil {
		return "", fmt.Errorf("get cookie: %w", err)
	}
	return token.Value, nil
}

func (c *cookieHelper) SetTokenInResponse(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     c.key,
		Value:    token,
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	})
}

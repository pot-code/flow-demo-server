package auth

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

type TokenData struct {
	SessionID string `json:"sid"`
}

func (t *TokenData) toClaim() jwt.Claims {
	return jwt.MapClaims{
		"sid": t.SessionID,
	}
}

func (t *TokenData) fromClaim(claims jwt.Claims) *TokenData {
	c, ok := claims.(jwt.MapClaims)
	if !ok {
		panic("claims is not jwt.MapClaims")
	}

	t.SessionID = c["sid"].(string)
	return t
}

type TokenService interface {
	GenerateToken(user *TokenData) (string, error)
	Verify(token string) (*TokenData, error)
	WithHttpResponse(w http.ResponseWriter, token string)
	FromHttpRequest(r *http.Request) (string, error)
}

type jwtTokenService struct {
	secret string
	key    string
}

// FromHttpRequest implements TokenService.
func (s *jwtTokenService) FromHttpRequest(r *http.Request) (string, error) {
	token, err := r.Cookie(s.key)
	if err != nil {
		return "", fmt.Errorf("get cookie: %w", err)
	}
	return token.Value, nil
}

// WithHttpResponse implements TokenService.
func (s *jwtTokenService) WithHttpResponse(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     s.key,
		Value:    token,
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	})
}

func (s *jwtTokenService) GenerateToken(u *TokenData) (string, error) {
	return s.Sign(u.toClaim())
}

func (s *jwtTokenService) Verify(token string) (*TokenData, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("parse token: %w", err)
	}
	return new(TokenData).fromClaim(t.Claims), nil
}

func (s *jwtTokenService) Sign(claims jwt.Claims) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(s.secret))
}

func NewJwtTokenService(secret string, key string) TokenService {
	return &jwtTokenService{secret: secret, key: key}
}

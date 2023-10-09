package token

import (
	"fmt"

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

type Service interface {
	GenerateToken(user *TokenData) (string, error)
	Verify(token string) (*TokenData, error)
}

type service struct {
	secret string
}

func (s *service) GenerateToken(u *TokenData) (string, error) {
	return s.Sign(u.toClaim())
}

func (s *service) Verify(token string) (*TokenData, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("parse token: %w", err)
	}
	return new(TokenData).fromClaim(t.Claims), nil
}

func (s *service) Sign(claims jwt.Claims) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(s.secret))
}

func NewService(secret string) *service {
	return &service{secret: secret}
}

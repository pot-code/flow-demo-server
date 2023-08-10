package auth

import (
	"context"
	"gobit-demo/internal/token"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenService interface {
	GenerateToken(user *LoginUser) (string, error)
	Verify(token string) (*LoginUser, error)
	AddToBlacklist(ctx context.Context, token string) error
	IsInBlacklist(ctx context.Context, token string) (bool, error)
}

type jwtService struct {
	i   *token.JwtIssuer
	bl  TokenBlacklist
	exp time.Duration
}

func NewJwtTokenService(bl TokenBlacklist, secret string, tokenExpiration time.Duration) TokenService {
	return &jwtService{i: token.NewJwtIssuer(secret), bl: bl, exp: tokenExpiration}
}

func (s *jwtService) GenerateToken(u *LoginUser) (string, error) {
	return s.i.Sign(u.toClaim(s.exp))
}

func (s *jwtService) Verify(token string) (*LoginUser, error) {
	c, err := s.i.Verify(token)
	if err != nil {
		return nil, err
	}
	return new(LoginUser).fromClaim(c), nil
}

func (s *jwtService) AddToBlacklist(ctx context.Context, token string) error {
	return s.bl.Add(ctx, token)
}

func (s *jwtService) IsInBlacklist(ctx context.Context, token string) (bool, error) {
	return s.bl.Has(ctx, token)
}

func (u *LoginUser) toClaim(exp time.Duration) jwt.Claims {
	return jwt.MapClaims{
		"id":       u.Id,
		"username": u.Username,
		"name":     u.Name,
		"exp":      float64(time.Now().Add(exp).Unix()),
	}
}

func (u *LoginUser) fromClaim(claims jwt.Claims) *LoginUser {
	c, ok := claims.(jwt.MapClaims)
	if !ok {
		panic("claims is not jwt.MapClaims")
	}

	u.Id = uint(c["id"].(float64))
	u.Username = c["username"].(string)
	u.Name = c["name"].(string)
	return u
}

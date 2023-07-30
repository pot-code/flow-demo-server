package auth

import (
	"context"
	"errors"
	"fmt"
	"gobit-demo/ent"
	"gobit-demo/ent/user"
	"gobit-demo/internal/db"
	"gobit-demo/internal/token"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound         = errors.New("用户未找到")
	ErrDuplicatedUser       = errors.New("用户已存在")
	ErrIncorrectCredentials = errors.New("用户名或密码错误")
)

type AuthService struct {
	e *ent.Client
}

func NewAuthService(client *ent.Client) *AuthService {
	return &AuthService{e: client}
}

func (s *AuthService) FindUserByUserName(ctx context.Context, name string) (*LoginUser, error) {
	user, err := s.e.User.Query().
		Where(user.Name(name)).Only(ctx)
	if ent.IsNotFound(err) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find user by name: %w", err)
	}
	return new(LoginUser).fromUser(user), err
}

func (s *AuthService) FindUserByMobile(ctx context.Context, mobile string) (*LoginUser, error) {
	user, err := s.e.User.Query().
		Where(user.Mobile(mobile)).Only(ctx)
	if ent.IsNotFound(err) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find user by mobile: %w", err)
	}
	return new(LoginUser).fromUser(user), err
}

func (s *AuthService) CreateUser(ctx context.Context, payload *createUserRequest) error {
	return db.WithEntTx(ctx, s.e, func(tx *ent.Tx) error {
		ok, err := tx.User.Query().
			Where(user.Or(user.Username(payload.Username), user.Mobile(payload.Mobile))).
			Exist(ctx)
		if err != nil {
			return fmt.Errorf("check duplicate user: %w", err)
		}
		if ok {
			return ErrDuplicatedUser
		}

		h, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("hash password: %w", err)
		}

		if _, err = tx.User.Create().
			SetUsername(payload.Username).
			SetName(payload.Name).
			SetMobile(payload.Mobile).
			SetPassword(string(h)).
			Save(ctx); err != nil {
			return fmt.Errorf("create user: %w", err)
		}

		return nil
	})
}

func (s *AuthService) FindUserByCredential(ctx context.Context, req *loginRequest) (*LoginUser, error) {
	user, err := s.e.User.Query().
		Where(user.Or(user.Username(req.Username), user.Mobile(req.Mobile))).
		Only(ctx)
	if ent.IsNotFound(err) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find user by login credentials: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, ErrIncorrectCredentials
	}

	return new(LoginUser).fromUser(user), err
}

type JwtService struct {
	jwt *token.JwtIssuer
	exp time.Duration
}

func NewJwtService(jwt *token.JwtIssuer, exp time.Duration) *JwtService {
	return &JwtService{jwt: jwt}
}

func (s *JwtService) GenerateToken(u *LoginUser) (string, error) {
	return s.jwt.Sign(s.userToClaim(u))
}

func (j *JwtService) userToClaim(u *LoginUser) jwt.Claims {
	return jwt.MapClaims{
		"id":       u.Id,
		"username": u.Username,
		"name":     u.Name,
		"mobile":   u.Mobile,
		"exp":      float64(time.Now().Add(j.exp).Unix()),
	}
}

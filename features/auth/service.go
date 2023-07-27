package auth

import (
	"context"
	"errors"
	"fmt"
	"gobit-demo/ent"
	"gobit-demo/ent/user"
	"gobit-demo/internal/token"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound         = errors.New("用户未找到")
	ErrDuplicatedUser       = errors.New("用户已存在")
	ErrIncorrectCredentials = errors.New("用户名或密码错误")
)

type Service struct {
	jwt *token.JwtIssuer
	e   *ent.Client
	exp time.Duration
}

func NewService(client *ent.Client, jwt *token.JwtIssuer, exp time.Duration) *Service {
	return &Service{e: client, jwt: jwt, exp: exp}
}

func (s *Service) FindUserByUserName(ctx context.Context, name string) (*LoginUser, error) {
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

func (s *Service) FindUserByMobile(ctx context.Context, mobile string) (*LoginUser, error) {
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

func (s *Service) CreateUser(ctx context.Context, dto *CreateUserRequest) error {
	tx, err := s.e.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Commit()

	ok, err := tx.User.Query().
		Where(user.Or(user.Username(dto.Username), user.Mobile(dto.Mobile))).
		Exist(ctx)
	if err != nil {
		return fmt.Errorf("check duplicate user: %w", err)
	}
	if ok {
		return ErrDuplicatedUser
	}

	h, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	if _, err = tx.User.Create().
		SetUsername(dto.Username).
		SetName(dto.Name).
		SetMobile(dto.Mobile).
		SetPassword(string(h)).
		Save(ctx); err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	return nil
}

func (s *Service) FindUserByCredential(ctx context.Context, req *LoginRequest) (*LoginUser, error) {
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

func (s *Service) CreateToken(ctx context.Context, user *LoginUser) (string, error) {
	return s.jwt.Sign(user.toClaim(float64(time.Now().Add(s.exp).Unix())))
}

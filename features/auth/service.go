package auth

import (
	"context"
	"errors"
	"fmt"
	"gobit-demo/internal/token"
	"gobit-demo/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound         = errors.New("用户未找到")
	ErrDuplicatedUser       = errors.New("用户已存在")
	ErrIncorrectCredentials = errors.New("用户名或密码错误")
)

type AuthService struct {
	g *gorm.DB
}

func NewAuthService(g *gorm.DB) *AuthService {
	return &AuthService{g: g}
}

func (s *AuthService) FindUserByUserName(ctx context.Context, username string) (*LoginUser, error) {
	user := new(model.User)
	err := s.g.WithContext(ctx).Where(&model.User{Username: username}).First(user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find user by name: %w", err)
	}
	return new(LoginUser).fromUser(user), err
}

func (s *AuthService) FindUserByMobile(ctx context.Context, mobile string) (*LoginUser, error) {
	user := new(model.User)
	err := s.g.WithContext(ctx).Where(&model.User{Mobile: mobile}).First(user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find user by mobile: %w", err)
	}
	return new(LoginUser).fromUser(user), err
}

func (s *AuthService) CreateUser(ctx context.Context, payload *createUserRequest) error {
	return s.g.Transaction(func(tx *gorm.DB) error {
		var count int64
		if err := s.g.WithContext(ctx).
			Model(&model.User{}).
			Where(&model.User{Mobile: payload.Mobile}).
			Or(&model.User{Username: payload.Username}).
			Count(&count).
			Error; err != nil {
			return fmt.Errorf("check duplicate user: %w", err)
		}
		if count > 0 {
			return ErrDuplicatedUser
		}

		h, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("hash password: %w", err)
		}

		if err = s.g.Create(&model.User{
			Name:     payload.Name,
			Username: payload.Username,
			Password: string(h),
			Mobile:   payload.Mobile,
		}).Error; err != nil {
			return fmt.Errorf("create user: %w", err)
		}

		return nil
	})
}

func (s *AuthService) FindUserByCredential(ctx context.Context, req *loginRequest) (*LoginUser, error) {
	user := new(model.User)
	err := s.g.WithContext(ctx).
		Where(&model.User{Mobile: req.Mobile}).
		Or(&model.User{Username: req.Username}).
		First(user).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
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
	return &JwtService{jwt: jwt, exp: exp}
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

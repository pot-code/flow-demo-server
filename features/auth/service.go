package auth

import (
	"context"
	"errors"
	"fmt"
	"gobit-demo/internal/event"
	"gobit-demo/internal/token"
	"gobit-demo/internal/util"
	"gobit-demo/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound         = errors.New("用户不存在")
	ErrDuplicatedUser       = errors.New("用户已存在")
	ErrIncorrectCredentials = errors.New("用户名或密码错误")
)

type AuthService struct {
	g  *gorm.DB
	eb event.EventBus
}

func NewAuthService(g *gorm.DB, eb event.EventBus) *AuthService {
	return &AuthService{g: g, eb: eb}
}

func (s *AuthService) FindUserByUserName(ctx context.Context, username string) (*LoginUser, error) {
	user := new(LoginUser)
	err := s.g.WithContext(ctx).Model(&model.User{}).Where(&model.User{Username: username}).Take(user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find user by name: %w", err)
	}
	return user, err
}

func (s *AuthService) FindUserByMobile(ctx context.Context, mobile string) (*LoginUser, error) {
	user := new(LoginUser)
	err := s.g.WithContext(ctx).Model(&model.User{}).Where(&model.User{Mobile: mobile}).Take(user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find user by mobile: %w", err)
	}
	return user, err
}

func (s *AuthService) CreateUser(ctx context.Context, payload *CreateUserRequest) (*RegisterUser, error) {
	user := model.User{
		Name:     payload.Name,
		Username: payload.Username,
		Mobile:   payload.Mobile,
	}
	if err := s.g.Transaction(func(tx *gorm.DB) error {
		ok, err := util.GormCheckExistence(s.g, func(tx *gorm.DB) *gorm.DB {
			return tx.WithContext(ctx).
				Model(&model.User{}).
				Select("1").
				Where(&model.User{Mobile: payload.Mobile}).
				Or(&model.User{Username: payload.Username}).Take(nil)
		})
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
		user.Password = string(h)

		if err = tx.WithContext(ctx).Create(&user).Error; err != nil {
			return fmt.Errorf("create user: %w", err)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	s.eb.Publish(&UserCreatedEvent{
		ID:       user.ID,
		Username: user.Username,
	})
	return new(RegisterUser).fromUser(&user), nil
}

func (s *AuthService) FindUserByCredential(ctx context.Context, req *LoginRequest) (*LoginUser, error) {
	user := new(model.User)
	err := s.g.WithContext(ctx).
		Where(&model.User{Mobile: req.Mobile}).
		Or(&model.User{Username: req.Username}).
		Take(user).
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

	s.eb.Publish(&UserLoginEvent{
		ID:       user.ID,
		Username: user.Username,
	})
	return new(LoginUser).fromUser(user), err
}

func (u *LoginUser) fromUser(user *model.User) *LoginUser {
	u.Id = user.ID
	u.Name = user.Name
	u.Username = user.Username
	u.Mobile = user.Mobile
	return u
}

func (u *RegisterUser) fromUser(user *model.User) *RegisterUser {
	u.Id = user.ID
	u.Name = user.Name
	u.Username = user.Username
	u.Mobile = user.Mobile
	return u
}

type tokenBlacklist interface {
	Add(ctx context.Context, token string) error
	Has(ctx context.Context, token string) (bool, error)
}

type JwtService struct {
	jwt *token.JwtIssuer
	bl  tokenBlacklist
	exp time.Duration
}

func NewJwtService(jwt *token.JwtIssuer, bl tokenBlacklist, exp time.Duration) *JwtService {
	return &JwtService{jwt: jwt, bl: bl, exp: exp}
}

func (s *JwtService) GenerateToken(u *LoginUser) (string, error) {
	return s.jwt.Sign(u.toClaim(s.exp))
}

func (s *JwtService) Verify(token string) (jwt.Claims, error) {
	return s.jwt.Verify(token)
}

func (s *JwtService) AddToBlacklist(ctx context.Context, token string) error {
	return s.bl.Add(ctx, token)
}

func (s *JwtService) IsInBlacklist(ctx context.Context, token string) (bool, error) {
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

type redisTokenBlacklist struct {
	rc  *redis.Client
	exp time.Duration
}

func NewRedisTokenBlacklist(rc *redis.Client, exp time.Duration) *redisTokenBlacklist {
	return &redisTokenBlacklist{rc: rc, exp: exp}
}

func (r *redisTokenBlacklist) Add(ctx context.Context, token string) error {
	if err := r.rc.Set(ctx, r.getKey(token), 1, r.exp).Err(); err != nil {
		return fmt.Errorf("add token to blacklist: %w", err)
	}
	return nil
}

func (r *redisTokenBlacklist) Has(ctx context.Context, token string) (bool, error) {
	c, err := r.rc.Exists(ctx, r.getKey(token)).Result()
	if err != nil {
		return false, fmt.Errorf("check token in blacklist: %w", err)
	}
	return c == 1, nil
}

func (r *redisTokenBlacklist) getKey(token string) string {
	return fmt.Sprintf("jwt:blacklist:%s", token)
}

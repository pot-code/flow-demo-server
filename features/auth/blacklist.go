package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type TokenBlacklist interface {
	Add(ctx context.Context, token string) error
	Has(ctx context.Context, token string) (bool, error)
	Delete(ctx context.Context, token string) (bool, error)
}

type redisTokenBlacklist struct {
	rc  *redis.Client
	exp time.Duration
}

func NewRedisTokenBlacklist(rc *redis.Client, expiration time.Duration) TokenBlacklist {
	return &redisTokenBlacklist{rc: rc, exp: expiration}
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

func (r *redisTokenBlacklist) Delete(ctx context.Context, token string) (bool, error) {
	c, err := r.rc.Del(ctx, r.getKey(token)).Result()
	if err != nil {
		return false, fmt.Errorf("delete token from blacklist: %w", err)
	}
	return c == 1, nil
}

func (r *redisTokenBlacklist) getKey(token string) string {
	return fmt.Sprintf("token:blacklist:%s", token)
}

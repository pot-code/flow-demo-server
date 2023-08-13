package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type TokenBlockList interface {
	Add(ctx context.Context, token string) error
	Has(ctx context.Context, token string) (bool, error)
	Delete(ctx context.Context, token string) (bool, error)
}

type redisTokenBlockList struct {
	rc  *redis.Client
	exp time.Duration
}

func NewRedisTokenBlacklist(rc *redis.Client, expiration time.Duration) TokenBlockList {
	return &redisTokenBlockList{rc: rc, exp: expiration}
}

func (r *redisTokenBlockList) Add(ctx context.Context, token string) error {
	if err := r.rc.Set(ctx, r.getKey(token), 1, r.exp).Err(); err != nil {
		return fmt.Errorf("add token to block list: %w", err)
	}
	return nil
}

func (r *redisTokenBlockList) Has(ctx context.Context, token string) (bool, error) {
	c, err := r.rc.Exists(ctx, r.getKey(token)).Result()
	if err != nil {
		return false, fmt.Errorf("check token in block list: %w", err)
	}
	return c == 1, nil
}

func (r *redisTokenBlockList) Delete(ctx context.Context, token string) (bool, error) {
	c, err := r.rc.Del(ctx, r.getKey(token)).Result()
	if err != nil {
		return false, fmt.Errorf("delete token from block list: %w", err)
	}
	return c == 1, nil
}

func (r *redisTokenBlockList) getKey(token string) string {
	return fmt.Sprintf("token:block:%s", token)
}

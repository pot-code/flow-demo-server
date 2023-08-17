package auth

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var (
	ErrSessionNotFound = fmt.Errorf("session not found")
)

type Session struct {
	SessionID       string
	UserID          uint
	Username        string
	UserPermissions []string
	UserRoles       []string
}

type SessionManager interface {
	GetSession(ctx context.Context, sid string) (*Session, error)
	NewSession(ctx context.Context, uid uint, username string, permissions []string, roles []string) (*Session, error)
	DeleteSession(ctx context.Context, sid string) error
}

type redisSessionManager struct {
	r   *redis.Client
	exp time.Duration
}

// DeleteSession implements SessionManager.
func (s *redisSessionManager) DeleteSession(ctx context.Context, sid string) error {
	key := s.getRedisKey(sid)
	if err := s.r.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("delete session from redis: %w", err)
	}
	return nil
}

// GetSession implements SessionManager.
func (s *redisSessionManager) GetSession(ctx context.Context, sid string) (*Session, error) {
	key := s.getRedisKey(sid)

	code, err := s.r.Exists(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("check session in redis: %w", err)
	}
	if code == 0 {
		return nil, ErrSessionNotFound
	}

	data, err := s.r.Get(ctx, key).Bytes()
	if err != nil {
		return nil, fmt.Errorf("get session from redis: %w", err)
	}

	session := new(Session)
	if err := gob.NewDecoder(bytes.NewBuffer(data)).Decode(session); err != nil {
		return nil, fmt.Errorf("decode session: %w", err)
	}
	return session, nil
}

// NewSession implements SessionManager.
func (s *redisSessionManager) NewSession(
	ctx context.Context,
	uid uint,
	username string,
	permissions []string,
	roles []string,
) (*Session, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, fmt.Errorf("generate uuid: %w", err)
	}

	session := &Session{
		SessionID:       id.String(),
		UserID:          uid,
		Username:        username,
		UserPermissions: permissions,
		UserRoles:       roles,
	}
	key := s.getRedisKey(id.String())
	bs := new(bytes.Buffer)
	if err := gob.NewEncoder(bs).Encode(session); err != nil {
		return nil, fmt.Errorf("encode session: %w", err)
	}
	if err := s.r.Set(ctx, key, bs.Bytes(), s.exp).Err(); err != nil {
		return nil, fmt.Errorf("set session to redis: %w", err)
	}
	return session, nil
}

func (s *redisSessionManager) getRedisKey(sid string) string {
	return fmt.Sprintf("auth:session:%s", sid)
}

func NewRedisSessionManager(r *redis.Client, expiration time.Duration) SessionManager {
	return &redisSessionManager{r: r, exp: expiration}
}

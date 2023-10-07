package auth

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"gobit-demo/model"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type sessionKeyType struct{}

var sessionKey = sessionKeyType{}

var (
	ErrSessionNotFound = fmt.Errorf("session not found")
)

type Session struct {
	SessionID       string
	UserID          model.ID
	Username        string
	UserPermissions []string
	UserRoles       []string
	ExpiredAt       time.Time
}

type SessionManager interface {
	GetSessionBySessionID(ctx context.Context, sid string) (*Session, error)
	GetSessionFromContext(ctx context.Context) *Session
	SetSession(ctx context.Context, s *Session) context.Context
	NewSession(ctx context.Context, uid model.ID, username string, permissions []string, roles []string) (*Session, error)
	DeleteSession(ctx context.Context, sid string) error
	RefreshSession(ctx context.Context, s *Session, threshold time.Duration) error
}

type redisSessionManager struct {
	r   *redis.Client
	exp time.Duration
}

// RefreshSession implements SessionManager.
func (sm *redisSessionManager) RefreshSession(ctx context.Context, s *Session, threshold time.Duration) error {
	ttl := time.Until(s.ExpiredAt)
	if ttl > threshold {
		return nil
	}

	key := sm.getRedisKey(s.SessionID)
	if err := sm.r.Expire(ctx, key, sm.exp).Err(); err != nil {
		return fmt.Errorf("refresh session: %w", err)
	}
	return nil
}

// GetSessionFromContext implements SessionManager.
func (sm *redisSessionManager) GetSessionFromContext(ctx context.Context) *Session {
	v, ok := ctx.Value(sessionKey).(*Session)
	if !ok {
		panic("session not found in context")
	}
	return v
}

// SetSession implements SessionManager.
func (sm *redisSessionManager) SetSession(ctx context.Context, s *Session) context.Context {
	return context.WithValue(ctx, sessionKey, s)
}

// DeleteSession implements SessionManager.
func (sm *redisSessionManager) DeleteSession(ctx context.Context, sid string) error {
	key := sm.getRedisKey(sid)
	if err := sm.r.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("delete session from redis: %w", err)
	}
	log.Trace().Str("sid", sid).Msg("session deleted")
	return nil
}

// GetSessionBySessionID implements SessionManager.
func (sm *redisSessionManager) GetSessionBySessionID(ctx context.Context, sid string) (*Session, error) {
	ok, err := sm.exists(ctx, sid)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrSessionNotFound
	}

	key := sm.getRedisKey(sid)
	data, err := sm.r.Get(ctx, key).Bytes()
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
func (sm *redisSessionManager) NewSession(
	ctx context.Context,
	uid model.ID,
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
		ExpiredAt:       time.Now().Add(sm.exp),
	}
	key := sm.getRedisKey(id.String())
	bs := new(bytes.Buffer)
	if err := gob.NewEncoder(bs).Encode(session); err != nil {
		return nil, fmt.Errorf("encode session: %w", err)
	}
	if err := sm.r.Set(ctx, key, bs.Bytes(), sm.exp).Err(); err != nil {
		return nil, fmt.Errorf("set session to redis: %w", err)
	}
	log.Trace().
		Str("sid", id.String()).
		Str("username", username).
		Strs("permissions", permissions).
		Strs("roles", roles).
		Msg("session created")
	return session, nil
}

func (sm *redisSessionManager) getRedisKey(sid string) string {
	return fmt.Sprintf("auth:session:%s", sid)
}

func (sm *redisSessionManager) exists(ctx context.Context, sid string) (bool, error) {
	key := sm.getRedisKey(sid)
	code, err := sm.r.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("check session in redis: %w", err)
	}
	if code == 0 {
		return false, nil
	}
	return true, nil
}

func NewRedisSessionManager(r *redis.Client, expiration time.Duration) SessionManager {
	return &redisSessionManager{r: r, exp: expiration}
}

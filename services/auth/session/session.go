package session

import (
	"context"
	"fmt"
	"gobit-demo/model"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

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
	NewSession(ctx context.Context, uid model.ID, username string, permissions []string, roles []string) (*Session, error)
	DeleteSession(ctx context.Context, sid string) error
	RefreshSession(ctx context.Context, s *Session, threshold time.Duration) error
}

type sessionManager struct {
	r   *redis.Client
	se  SessionSerializer
	exp time.Duration
}

// RefreshSession implements SessionManager.
func (sm *sessionManager) RefreshSession(ctx context.Context, s *Session, threshold time.Duration) error {
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

// DeleteSession implements SessionManager.
func (sm *sessionManager) DeleteSession(ctx context.Context, sid string) error {
	key := sm.getRedisKey(sid)
	if err := sm.r.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("delete session from redis: %w", err)
	}
	log.Trace().Str("sid", sid).Msg("session deleted")
	return nil
}

// GetSessionBySessionID implements SessionManager.
func (sm *sessionManager) GetSessionBySessionID(ctx context.Context, sid string) (*Session, error) {
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

	session, err := sm.se.Deserialize(data)
	if err != nil {
		return nil, fmt.Errorf("deserialize session: %w", err)
	}
	return session, nil
}

// NewSession implements SessionManager.
func (sm *sessionManager) NewSession(
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
	b, err := sm.se.Serialize(session)
	if err != nil {
		return nil, fmt.Errorf("serialize session: %w", err)
	}
	if err := sm.r.Set(ctx, key, b, sm.exp).Err(); err != nil {
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

func (sm *sessionManager) getRedisKey(sid string) string {
	return fmt.Sprintf("auth:session:%s", sid)
}

func (sm *sessionManager) exists(ctx context.Context, sid string) (bool, error) {
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

func NewSessionManager(r *redis.Client, se SessionSerializer, expiration time.Duration) *sessionManager {
	return &sessionManager{r: r, exp: expiration, se: se}
}

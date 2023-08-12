package audit

import (
	"context"
	"encoding/json"
	"fmt"
	"gobit-demo/features/auth"
	"gobit-demo/model"
	"reflect"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type AuditLog struct {
	a       *model.AuditLog
	g       *gorm.DB
	payload any
}

func NewAuditLog(g *gorm.DB) *AuditLog {
	return &AuditLog{a: new(model.AuditLog), g: g}
}

func (b *AuditLog) Subject(subject string) *AuditLog {
	b.a.Subject = subject
	return b
}

func (b *AuditLog) Action(action string) *AuditLog {
	b.a.Action = action
	return b
}

func (b *AuditLog) Payload(data any) *AuditLog {
	if reflect.TypeOf(data).Kind() != reflect.Pointer {
		panic("data must be pointer")
	}

	b.payload = data
	return b
}

func (b *AuditLog) Commit(ctx context.Context) error {
	if b.a.Action == "" && b.a.Subject == "" && b.payload == nil {
		log.Warn().Msg("empty audit log")
		return nil
	}

	if b.a.Subject == "" {
		u, ok := new(auth.LoginUser).FromContext(ctx)
		if ok {
			b.a.Subject = u.Username
		}
	}

	if b.payload != nil {
		bs, err := json.Marshal(b.payload)
		if err != nil {
			return fmt.Errorf("marshal data: %w", err)
		}
		b.a.Payload = string(bs)
	}

	if err := b.g.WithContext(ctx).Create(b.a).Error; err != nil {
		return fmt.Errorf("create audit log: %w", err)
	}
	return nil
}

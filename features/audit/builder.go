package audit

import (
	"context"
	"encoding/json"
	"fmt"
	"gobit-demo/features/auth"
	"reflect"

	"gorm.io/gorm"
)

type auditLogBuilder struct {
	a *auditLog
	g *gorm.DB
}

func newAuditLogBuilder(g *gorm.DB) *auditLogBuilder {
	return &auditLogBuilder{a: new(auditLog), g: g}
}

func (b *auditLogBuilder) Subject(subject string) *auditLogBuilder {
	b.a.Subject = subject
	return b
}

func (b *auditLogBuilder) Action(action string) *auditLogBuilder {
	b.a.Action = action
	return b
}

func (b *auditLogBuilder) Payload(data any) *auditLogBuilder {
	if reflect.TypeOf(data).Kind() != reflect.Pointer {
		panic("data must be pointer")
	}

	b.a.rawPayload = data
	return b
}

func (b *auditLogBuilder) Commit(ctx context.Context) error {
	if b.a.Subject == "" {
		u, ok := new(auth.LoginUser).FromContext(ctx)
		if ok {
			b.a.Subject = u.Username
		}
	}

	if b.a.rawPayload != nil {
		bs, err := json.Marshal(b.a.rawPayload)
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

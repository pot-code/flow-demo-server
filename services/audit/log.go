package audit

import (
	"context"
	"encoding/json"
	"fmt"
	"gobit-demo/model"
	"gobit-demo/services/auth/session"

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

func (a *AuditLog) Subject(subject string) *AuditLog {
	a.a.Subject = subject
	return a
}

func (a *AuditLog) Action(action string) *AuditLog {
	a.a.Action = action
	return a
}

func (a *AuditLog) Payload(data any) *AuditLog {
	a.payload = data
	return a
}

func (a *AuditLog) UseContext(ctx context.Context) *AuditLog {
	s := session.GetSessionFromContext(ctx)
	a.a.Subject = s.Username
	return a
}

func (a *AuditLog) Commit(ctx context.Context) error {
	if a.a.Action == "" && a.a.Subject == "" && a.payload == nil {
		panic("empty audit log")
	}
	if a.a.Subject == "" {
		panic("subject cannot be empty")
	}
	if a.payload != nil {
		bs, err := json.Marshal(a.payload)
		if err != nil {
			return fmt.Errorf("marshal data: %w", err)
		}
		a.a.Payload = string(bs)
	}

	if err := a.g.WithContext(ctx).Create(a.a).Error; err != nil {
		return fmt.Errorf("create audit log: %w", err)
	}
	return nil
}

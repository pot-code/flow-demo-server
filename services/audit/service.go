package audit

import (
	"gobit-demo/services/auth"

	"gorm.io/gorm"
)

type Service interface {
	NewAuditLog() *AuditLog
}

type service struct {
	g  *gorm.DB
	sm auth.SessionManager
}

func (s *service) NewAuditLog() *AuditLog {
	return NewAuditLog(s.g, s.sm)
}

func NewService(g *gorm.DB, sm auth.SessionManager) *service {
	return &service{g: g, sm: sm}
}

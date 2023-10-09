package audit

import (
	"gorm.io/gorm"
)

type Service interface {
	NewAuditLog() *AuditLog
}

type service struct {
	g *gorm.DB
}

func (s *service) NewAuditLog() *AuditLog {
	return NewAuditLog(s.g)
}

func NewService(g *gorm.DB) *service {
	return &service{g: g}
}

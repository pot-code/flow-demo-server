package audit

import (
	"gorm.io/gorm"
)

type Service interface {
	NewAuditLog() *AuditLog
}

func NewService(g *gorm.DB) Service {
	return &service{g: g}
}

type service struct {
	g *gorm.DB
}

func (s *service) NewAuditLog() *AuditLog {
	return NewAuditLog(s.g)
}

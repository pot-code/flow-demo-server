package model

import "gorm.io/gorm"

type AuditLog struct {
	gorm.Model
	Subject string
	Action  string
	Payload string
}

package model

import (
	"time"
)

type AuditLog struct {
	ID        uint64 `gorm:"primaryKey,autoIncrement" json:"id,omitempty"`
	Subject   string
	Action    string
	Payload   string
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at,omitempty"`
}

package model

import (
	"gobit-demo/model/pk"
	"time"

	"github.com/segmentio/ksuid"
	"gorm.io/gorm"
)

type Notification struct {
	ID        pk.ID     `json:"id,omitempty"`
	Content   string    `json:"content,omitempty"`
	IsRead    bool      `json:"is_read,omitempty"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at,omitempty"`
	OwnerID   *pk.ID    `json:"owner_id,omitempty"`
	Owner     *User     `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
}

func (n *Notification) BeforeCreate(tx *gorm.DB) error {
	n.ID = pk.ParseFromString(ksuid.New().String())
	return nil
}

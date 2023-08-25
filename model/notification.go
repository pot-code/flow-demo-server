package model

import (
	"gobit-demo/internal/uuid"
	"time"

	"gorm.io/gorm"
)

type Notification struct {
	ID        UUID      `json:"id,omitempty"`
	Content   string    `json:"content,omitempty"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at,omitempty"`
	OwnerID   *UUID     `json:"owner_id,omitempty"`
	Owner     *User     `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
}

func (n *Notification) BeforeCreate(tx *gorm.DB) error {
	uid, err := uuid.Snoyflake.NextID()
	if err != nil {
		return err
	}
	n.ID = UUID(uid)
	return nil
}

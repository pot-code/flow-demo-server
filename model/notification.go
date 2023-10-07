package model

import (
	"gobit-demo/infra/uuid"
	"time"

	"gorm.io/gorm"
)

type Notification struct {
	ID        ID        `json:"id,omitempty"`
	Content   string    `json:"content,omitempty"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at,omitempty"`
	OwnerID   *ID       `json:"owner_id,omitempty"`
	Owner     *User     `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
}

func (n *Notification) BeforeCreate(tx *gorm.DB) error {
	uid, err := uuid.Sonyflake.NextID()
	if err != nil {
		return err
	}
	n.ID = ID(uid)
	return nil
}

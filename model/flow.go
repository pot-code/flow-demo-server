package model

import (
	"gobit-demo/infra/uuid"
	"time"

	"gorm.io/gorm"
)

type Flow struct {
	ID          ID             `gorm:"primaryKey;type:BIGINT UNSIGNED" json:"id,omitempty"`
	Name        string         `gorm:"index,not null,size:32" json:"name,omitempty"`
	Description string         `json:"description,omitempty"`
	Nodes       string         `json:"nodes,omitempty"`
	Edges       string         `json:"edges,omitempty"`
	OwnerID     *ID            `json:"owner_id,omitempty"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Owner       *User          `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
}

func (f *Flow) BeforeCreate(tx *gorm.DB) error {
	uid, err := uuid.Sonyflake.NextID()
	if err != nil {
		return err
	}
	f.ID = ID(uid)
	return nil
}

package model

import (
	"gobit-demo/internal/uuid"

	"gorm.io/gorm"
)

type Flow struct {
	ID          UUID   `gorm:"primaryKey;type:BIGINT UNSIGNED" json:"id,omitempty"`
	Name        string `gorm:"index,not null,size:32" json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Nodes       string `json:"nodes,omitempty"`
	Edges       string `json:"edges,omitempty"`
	OwnerID     *UUID  `json:"owner_id,omitempty"`
	Owner       *User  `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
}

func (f *Flow) BeforeCreate(tx *gorm.DB) error {
	uid, err := uuid.Snoyflake.NextID()
	if err != nil {
		return err
	}
	f.ID = UUID(uid)
	return nil
}

package model

import (
	"gobit-demo/internal/uuid"

	"gorm.io/gorm"
)

type Flow struct {
	gorm.Model
	ID          UUID   `gorm:"primaryKey;type:bigint"`
	Name        string `gorm:"index,not null,size:32"`
	Description string
	Nodes       string
	Edges       string
	OwnerID     *UUID
	Owner       *User `gorm:"foreignKey:OwnerID"`
}

func (f *Flow) BeforeCreate(tx *gorm.DB) error {
	uid, err := uuid.Snoyflake.NextID()
	if err != nil {
		return err
	}
	f.ID = UUID(uid)
	return nil
}

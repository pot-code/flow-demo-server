package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Flow struct {
	gorm.Model
	ID          uuid.UUID `gorm:"primaryKey;type:char(36)"`
	Name        string    `gorm:"index,not null,size:32"`
	Description string
	Nodes       string
	Edges       string
	OwnerID     *string
	Owner       *User `gorm:"foreignKey:OwnerID"`
}

func (f *Flow) BeforeCreate(tx *gorm.DB) error {
	uid, err := uuid.NewUUID()
	if err != nil {
		return err
	}
	f.ID = uid
	return nil
}

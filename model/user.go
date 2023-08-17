package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       uuid.UUID `gorm:"primaryKey;type:char(36)"`
	Name     string
	Username string `gorm:"uniqueIndex,not null,size:12"`
	Password string `gorm:"not null"`
	Mobile   string `gorm:"uniqueIndex,not null,size:11"`
	Disabled bool
	Roles    []*Role `gorm:"many2many:user_roles"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	uid, err := uuid.NewUUID()
	if err != nil {
		return err
	}
	u.ID = uid
	return nil
}

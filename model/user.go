package model

import (
	"gobit-demo/internal/uuid"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       UUID `gorm:"primaryKey;type:bigint"`
	Name     string
	Username string `gorm:"uniqueIndex,not null,size:12"`
	Password string `gorm:"not null"`
	Mobile   string `gorm:"uniqueIndex,not null,size:11"`
	Disabled bool
	Roles    []*Role `gorm:"many2many:user_roles"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	uid, err := uuid.Snoyflake.NextID()
	if err != nil {
		return err
	}
	u.ID = UUID(uid)
	return nil
}

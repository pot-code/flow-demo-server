package model

import (
	"gobit-demo/internal/uuid"

	"gorm.io/gorm"
)

type User struct {
	ID       UUID    `gorm:"primaryKey;type:BIGINT UNSIGNED" json:"id,omitempty"`
	Name     string  `json:"name,omitempty"`
	Username string  `gorm:"uniqueIndex,not null,size:12" json:"username,omitempty"`
	Password string  `gorm:"not null" json:"password,omitempty"`
	Mobile   string  `gorm:"uniqueIndex,not null,size:11" json:"mobile,omitempty"`
	Disabled bool    `json:"disabled,omitempty"`
	Roles    []*Role `gorm:"many2many:user_roles" json:"roles,omitempty"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	uid, err := uuid.Snoyflake.NextID()
	if err != nil {
		return err
	}
	u.ID = UUID(uid)
	return nil
}

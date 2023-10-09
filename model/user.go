package model

import (
	"time"

	"github.com/segmentio/ksuid"
	"gorm.io/gorm"
)

type User struct {
	ID        ID             `gorm:"primaryKey;type:varchar(27)" json:"id,omitempty"`
	Name      string         `json:"name,omitempty"`
	Username  string         `gorm:"uniqueIndex,not null,size:12" json:"username,omitempty"`
	Password  string         `gorm:"not null" json:"password,omitempty"`
	Mobile    string         `gorm:"uniqueIndex,not null,size:11" json:"mobile,omitempty"`
	Disabled  bool           `json:"disabled,omitempty"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Roles     []*Role        `gorm:"many2many:user_roles" json:"roles,omitempty"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.ID = ID(ksuid.New().String())
	return nil
}

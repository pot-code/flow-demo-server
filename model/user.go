package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string
	Username string  `gorm:"uniqueIndex,not null,size:12"`
	Password string  `gorm:"not null"`
	Mobile   string  `gorm:"uniqueIndex,not null,size:11"`
	Roles    []*Role `gorm:"many2many:user_roles"`
}

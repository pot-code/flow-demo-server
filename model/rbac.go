package model

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	Name        string `gorm:"uniqueIndex,not null,size:24"`
	Description string
	Users       []*User       `gorm:"many2many:user_roles"`
	Permissions []*Permission `gorm:"many2many:role_permissions"`
}

type Permission struct {
	gorm.Model
	Object string
	Action string
}

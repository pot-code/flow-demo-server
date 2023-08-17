package model

import "gorm.io/gorm"

type Flow struct {
	gorm.Model
	Name        string `gorm:"uniqueIndex,not null,size:12"`
	Description string
	Nodes       string
	Edges       string
	OwnerID     *uint
	Owner       *User `gorm:"foreignKey:OwnerID"`
}

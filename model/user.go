package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string
	Username string `gorm:"uniqueIndex,size:12"`
	Password string
	Mobile   string `gorm:"uniqueIndex,size:11"`
}

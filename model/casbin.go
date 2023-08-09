package model

import "gorm.io/gorm"

type CasbinRule struct {
	gorm.Model
	Ptype string `gorm:"size:512;uniqueIndex"`
	V0    string `gorm:"size:512;uniqueIndex"`
	V1    string `gorm:"size:512;uniqueIndex"`
	V2    string `gorm:"size:512;uniqueIndex"`
	V3    string `gorm:"size:512;uniqueIndex"`
	V4    string `gorm:"size:512;uniqueIndex"`
	V5    string `gorm:"size:512;uniqueIndex"`
}

package model

import "gorm.io/gorm"

type Flow struct {
	gorm.Model
	Name        string `gorm:"uniqueIndex,not null,size:12"`
	Description string
	OwnerID     *uint
	Nodes       []*FlowNode
	Owner       *User `gorm:"foreignKey:OwnerID"`
}

type FlowNode struct {
	gorm.Model
	Name        string `gorm:"not null"`
	Description string
	FlowID      uint
	NextID      *uint
	PrevID      *uint
	Next        *FlowNode `gorm:"foreignKey:NextID"`
	Prev        *FlowNode `gorm:"foreignKey:PrevID"`
}

package model

import "gorm.io/gorm"

type Flow struct {
	gorm.Model
	Name        string `gorm:"uniqueIndex,size:12"`
	Description string
	Nodes       []*FlowNode
}

type FlowNode struct {
	gorm.Model
	Name        string
	Description string
	FlowID      uint
	NextID      uint
	PrevID      uint
	Next        *FlowNode `gorm:"foreignKey:NextID"`
	Prev        *FlowNode `gorm:"foreignKey:PrevID"`
}

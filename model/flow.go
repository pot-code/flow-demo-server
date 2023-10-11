package model

import (
	"gobit-demo/model/pk"
	"time"

	"github.com/segmentio/ksuid"
	"gorm.io/gorm"
)

type Flow struct {
	ID          pk.ID          `gorm:"primaryKey;type:varchar(27)" json:"id,omitempty"`
	Name        string         `gorm:"index,not null,size:32" json:"name,omitempty"`
	Description string         `json:"description,omitempty"`
	Nodes       string         `json:"nodes,omitempty"`
	Edges       string         `json:"edges,omitempty"`
	OwnerID     *pk.ID         `json:"owner_id,omitempty"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Owner       *User          `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
}

func (f *Flow) BeforeCreate(tx *gorm.DB) error {
	f.ID = pk.ParseFromString(ksuid.New().String())
	return nil
}

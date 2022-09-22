package db

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Group struct {
	ID        string `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	Users     []User     `gorm:"foreignKey:GroupID"`
	Resource  []Resource `gorm:"foreignKey:GroupID"`
}

func (group *Group) BeforeCreate(tx *gorm.DB) (err error) {
	// UUID version 4
	if group.ID == "" {
		group.ID = uuid.New().String()
	}
	return
}

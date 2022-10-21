package db

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Resource struct {
	ID        string `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Week      uint32
	Value     uint32
	UserID    string
	User      User // added this. Inverse pointer for the database
	GroupID   string
	Group     Group // added this. Inverse pointer for the database
}

func (r *Resource) BeforeCreate(tx *gorm.DB) (err error) {
	// UUID version 4
	if r.ID == "" {
		r.ID = uuid.New().String()
	}
	return
}

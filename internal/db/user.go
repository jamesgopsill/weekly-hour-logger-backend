package db

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID           string `gorm:"primaryKey"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Name         string
	Email        string `gorm:"uniqueIndex"`
	Scopes       Scopes
	PasswordHash string `json:"-"`
	GroupID      string
	Group        Group      // added this. Group ID above is the key, this is the inverse signpost?
	Resources    []Resource `gorm:"foreignKey:UserID"`
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	// UUID version 4
	if user.ID == "" {
		user.ID = uuid.New().String()
	}
	return
}

/*

func (b *Base) BeforeCreate(tx *gorm.DB) (err error) {
	// UUID version 4
	if b.ID == "" {
		b.ID = uuid.New().String()
	}
	return
}
*/

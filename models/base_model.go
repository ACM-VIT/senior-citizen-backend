package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Base struct {
	ID        string     `gorm:"primary_key;type:uuid;" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (b *Base) BeforeCreate(tx *gorm.DB) error {
	b.ID = uuid.New().String()
	return nil
}
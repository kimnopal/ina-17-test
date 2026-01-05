package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey;" json:"id"`
	Username string    `gorm:"unique;not null" json:"username"`
	Password string    `gorm:"not null" json:"password,omitempty"`
	gorm.Model
}

type UserResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

func (m *User) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}

package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Ticket struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	EventID  uuid.UUID `gorm:"type:uuid;not null" json:"event_id"`
	Category string    `gorm:"type:varchar(50)" json:"category"` // VIP, Regular
	Price    float64   `gorm:"type:decimal(12,2);not null" json:"price"`
	Quota    int       `gorm:"not null" json:"quota"`
	Event    Event     `gorm:"foreignKey:EventID" json:"event,omitempty"`
}

func (t *Ticket) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

type TicketResponse struct {
	ID       uuid.UUID `json:"id"`
	EventID  uuid.UUID `json:"event_id"`
	Category string    `json:"category"`
	Price    float64   `json:"price"`
	Quota    int       `json:"quota"`
}

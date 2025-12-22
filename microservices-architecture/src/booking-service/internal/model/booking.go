package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Booking struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	UserID      uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	EventID     uuid.UUID  `gorm:"type:uuid;not null" json:"event_id"`
	TicketID    uuid.UUID  `gorm:"type:uuid;not null" json:"ticket_id"`
	Quantity    int        `gorm:"not null" json:"quantity"`
	TotalAmount float64    `gorm:"type:decimal(12,2);not null" json:"total_amount"`
	Status      string     `gorm:"type:varchar(30);not null" json:"status"` // PENDING, CONFIRMED, CANCELLED
	ExpiredAt   *time.Time `gorm:"type:timestamp" json:"expired_at"`
	CreatedAt   time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	Event       Event      `gorm:"foreignKey:EventID" json:"event,omitempty"`
	Ticket      Ticket     `gorm:"foreignKey:TicketID" json:"ticket,omitempty"`
}

func (b *Booking) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}

type BookingResponse struct {
	ID          uuid.UUID  `json:"id"`
	UserID      uuid.UUID  `json:"user_id"`
	EventID     uuid.UUID  `json:"event_id"`
	TicketID    uuid.UUID  `json:"ticket_id"`
	Quantity    int        `json:"quantity"`
	TotalAmount float64    `json:"total_amount"`
	Status      string     `json:"status"`
	ExpiredAt   *time.Time `json:"expired_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

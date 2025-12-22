package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Payment struct {
	ID            uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	BookingID     uuid.UUID  `gorm:"type:uuid;not null" json:"booking_id"`
	UserID        uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	Amount        float64    `gorm:"type:decimal(12,2);not null" json:"amount"`
	Currency      string     `gorm:"type:varchar(10);default:IDR" json:"currency"`
	PaymentMethod string     `gorm:"type:varchar(50);not null" json:"payment_method"` // VA, EWALLET, QRIS
	Status        string     `gorm:"type:varchar(30);not null" json:"status"`         // PENDING, PAID, FAILED, EXPIRED
	ExpiredAt     *time.Time `gorm:"type:timestamp" json:"expired_at"`
	PaidAt        *time.Time `gorm:"type:timestamp" json:"paid_at"`
	CreatedAt     time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (p *Payment) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	if p.Currency == "" {
		p.Currency = "IDR"
	}
	return nil
}

type PaymentResponse struct {
	ID            uuid.UUID  `json:"id"`
	BookingID     uuid.UUID  `json:"booking_id"`
	UserID        uuid.UUID  `json:"user_id"`
	Amount        float64    `json:"amount"`
	Currency      string     `json:"currency"`
	PaymentMethod string     `json:"payment_method"`
	Status        string     `json:"status"`
	ExpiredAt     *time.Time `json:"expired_at,omitempty"`
	PaidAt        *time.Time `json:"paid_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

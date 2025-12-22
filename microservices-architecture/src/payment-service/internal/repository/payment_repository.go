package repository

import (
	"payment-service/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentRepository interface {
	Create(payment *model.Payment) error
	FindByID(id uuid.UUID) (*model.Payment, error)
	FindByBookingID(bookingID uuid.UUID) (*model.Payment, error)
	FindAll() ([]model.Payment, error)
	Update(payment *model.Payment) error
	UpdateStatus(id uuid.UUID, status string) error
}

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) Create(payment *model.Payment) error {
	return r.db.Create(payment).Error
}

func (r *paymentRepository) FindByID(id uuid.UUID) (*model.Payment, error) {
	var payment model.Payment
	err := r.db.First(&payment, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *paymentRepository) FindByBookingID(bookingID uuid.UUID) (*model.Payment, error) {
	var payment model.Payment
	err := r.db.First(&payment, "booking_id = ?", bookingID).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *paymentRepository) FindAll() ([]model.Payment, error) {
	var payments []model.Payment
	err := r.db.Find(&payments).Error
	return payments, err
}

func (r *paymentRepository) Update(payment *model.Payment) error {
	return r.db.Save(payment).Error
}

func (r *paymentRepository) UpdateStatus(id uuid.UUID, status string) error {
	return r.db.Model(&model.Payment{}).Where("id = ?", id).Update("status", status).Error
}

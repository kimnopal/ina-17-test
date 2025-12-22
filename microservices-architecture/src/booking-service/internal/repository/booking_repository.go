package repository

import (
	"booking-service/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BookingRepository interface {
	Create(booking *model.Booking) error
	FindByID(id uuid.UUID) (*model.Booking, error)
	FindByIDForUpdate(id uuid.UUID) (*model.Booking, error)
	FindAll() ([]model.Booking, error)
	FindByUserID(userID uuid.UUID) ([]model.Booking, error)
	Update(booking *model.Booking) error
	UpdateStatus(id uuid.UUID, status string) error
	WithTx(tx *gorm.DB) BookingRepository
}

type bookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) BookingRepository {
	return &bookingRepository{db: db}
}

func (r *bookingRepository) Create(booking *model.Booking) error {
	return r.db.Create(booking).Error
}

func (r *bookingRepository) FindByID(id uuid.UUID) (*model.Booking, error) {
	var booking model.Booking
	err := r.db.Preload("Event").Preload("Ticket").First(&booking, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &booking, nil
}

func (r *bookingRepository) FindByIDForUpdate(id uuid.UUID) (*model.Booking, error) {
	var booking model.Booking
	err := r.db.Clauses(clause.Locking{Strength: "UPDATE"}).
		Preload("Event").Preload("Ticket").
		First(&booking, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &booking, nil
}

func (r *bookingRepository) FindAll() ([]model.Booking, error) {
	var bookings []model.Booking
	err := r.db.Preload("Event").Preload("Ticket").Find(&bookings).Error
	return bookings, err
}

func (r *bookingRepository) FindByUserID(userID uuid.UUID) ([]model.Booking, error) {
	var bookings []model.Booking
	err := r.db.Preload("Event").Preload("Ticket").Where("user_id = ?", userID).Find(&bookings).Error
	return bookings, err
}

func (r *bookingRepository) Update(booking *model.Booking) error {
	return r.db.Save(booking).Error
}

func (r *bookingRepository) UpdateStatus(id uuid.UUID, status string) error {
	return r.db.Model(&model.Booking{}).Where("id = ?", id).Update("status", status).Error
}

func (r *bookingRepository) WithTx(tx *gorm.DB) BookingRepository {
	return &bookingRepository{db: tx}
}

package repository

import (
	"booking-service/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EventRepository interface {
	Create(event *model.Event) error
	FindByID(id uuid.UUID) (*model.Event, error)
	FindAll() ([]model.Event, error)
	Update(event *model.Event) error
	Delete(id uuid.UUID) error
}

type eventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) EventRepository {
	return &eventRepository{db: db}
}

func (r *eventRepository) Create(event *model.Event) error {
	return r.db.Create(event).Error
}

func (r *eventRepository) FindByID(id uuid.UUID) (*model.Event, error) {
	var event model.Event
	err := r.db.First(&event, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *eventRepository) FindAll() ([]model.Event, error) {
	var events []model.Event
	err := r.db.Find(&events).Error
	return events, err
}

func (r *eventRepository) Update(event *model.Event) error {
	return r.db.Save(event).Error
}

func (r *eventRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.Event{}, "id = ?", id).Error
}

package repository

import (
	"booking-service/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TicketRepository interface {
	Create(ticket *model.Ticket) error
	FindByID(id uuid.UUID) (*model.Ticket, error)
	FindByIDForUpdate(id uuid.UUID) (*model.Ticket, error)
	FindByEventID(eventID uuid.UUID) ([]model.Ticket, error)
	FindAll() ([]model.Ticket, error)
	Update(ticket *model.Ticket) error
	Delete(id uuid.UUID) error
	ReduceQuota(id uuid.UUID, quantity int) error
	IncreaseQuota(id uuid.UUID, quantity int) error
	WithTx(tx *gorm.DB) TicketRepository
}

type ticketRepository struct {
	db *gorm.DB
}

func NewTicketRepository(db *gorm.DB) TicketRepository {
	return &ticketRepository{db: db}
}

func (r *ticketRepository) Create(ticket *model.Ticket) error {
	return r.db.Create(ticket).Error
}

func (r *ticketRepository) FindByID(id uuid.UUID) (*model.Ticket, error) {
	var ticket model.Ticket
	err := r.db.Preload("Event").First(&ticket, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &ticket, nil
}

func (r *ticketRepository) FindByIDForUpdate(id uuid.UUID) (*model.Ticket, error) {
	var ticket model.Ticket
	err := r.db.Clauses(clause.Locking{Strength: "UPDATE"}).
		Preload("Event").
		First(&ticket, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &ticket, nil
}

func (r *ticketRepository) FindByEventID(eventID uuid.UUID) ([]model.Ticket, error) {
	var tickets []model.Ticket
	err := r.db.Where("event_id = ?", eventID).Find(&tickets).Error
	return tickets, err
}

func (r *ticketRepository) FindAll() ([]model.Ticket, error) {
	var tickets []model.Ticket
	err := r.db.Preload("Event").Find(&tickets).Error
	return tickets, err
}

func (r *ticketRepository) Update(ticket *model.Ticket) error {
	return r.db.Save(ticket).Error
}

func (r *ticketRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.Ticket{}, "id = ?", id).Error
}

func (r *ticketRepository) ReduceQuota(id uuid.UUID, quantity int) error {
	return r.db.Model(&model.Ticket{}).Where("id = ?", id).
		Update("quota", gorm.Expr("quota - ?", quantity)).Error
}

func (r *ticketRepository) IncreaseQuota(id uuid.UUID, quantity int) error {
	return r.db.Model(&model.Ticket{}).Where("id = ?", id).
		Update("quota", gorm.Expr("quota + ?", quantity)).Error
}

func (r *ticketRepository) WithTx(tx *gorm.DB) TicketRepository {
	return &ticketRepository{db: tx}
}

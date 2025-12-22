package service

import (
	"booking-service/internal/client"
	"booking-service/internal/model"
	"booking-service/internal/repository"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BookingService interface {
	CreateBooking(userID uuid.UUID, eventID uuid.UUID, ticketID uuid.UUID, quantity int) (*model.BookingResponse, error)
	GetBookingByID(id uuid.UUID) (*model.BookingResponse, error)
	GetAllBookings() ([]model.BookingResponse, error)
	UpdateBookingStatus(id uuid.UUID, status string) error
}

type bookingService struct {
	db            *gorm.DB
	bookingRepo   repository.BookingRepository
	ticketRepo    repository.TicketRepository
	eventRepo     repository.EventRepository
	userClient    client.UserClient
	paymentClient client.PaymentClient
	webhookClient client.WebhookClient
}

func NewBookingService(
	db *gorm.DB,
	bookingRepo repository.BookingRepository,
	ticketRepo repository.TicketRepository,
	eventRepo repository.EventRepository,
	userClient client.UserClient,
	paymentClient client.PaymentClient,
	webhookClient client.WebhookClient,
) BookingService {
	return &bookingService{
		db:            db,
		bookingRepo:   bookingRepo,
		ticketRepo:    ticketRepo,
		eventRepo:     eventRepo,
		userClient:    userClient,
		paymentClient: paymentClient,
		webhookClient: webhookClient,
	}
}

func (s *bookingService) CreateBooking(userID uuid.UUID, eventID uuid.UUID, ticketID uuid.UUID, quantity int) (*model.BookingResponse, error) {
	if quantity <= 0 {
		return nil, errors.New("quantity must be greater than 0")
	}

	event, err := s.eventRepo.FindByID(eventID)
	if err != nil {
		return nil, errors.New("event not found")
	}

	var booking *model.Booking
	var ticket *model.Ticket

	err = s.db.Transaction(func(tx *gorm.DB) error {
		ticketRepoTx := s.ticketRepo.WithTx(tx)
		bookingRepoTx := s.bookingRepo.WithTx(tx)

		ticket, err = ticketRepoTx.FindByIDForUpdate(ticketID)
		if err != nil {
			return errors.New("ticket not found")
		}

		if ticket.EventID != eventID {
			return errors.New("ticket does not belong to the specified event")
		}

		if ticket.Quota < quantity {
			return errors.New("insufficient ticket quota")
		}

		totalAmount := ticket.Price * float64(quantity)
		expiredAt := time.Now().Add(15 * time.Minute)

		booking = &model.Booking{
			UserID:      userID,
			EventID:     eventID,
			TicketID:    ticketID,
			Quantity:    quantity,
			TotalAmount: totalAmount,
			Status:      "PENDING",
			ExpiredAt:   &expiredAt,
		}

		if err := bookingRepoTx.Create(booking); err != nil {
			return err
		}

		if err := ticketRepoTx.ReduceQuota(ticketID, quantity); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	booking.Event = *event
	booking.Ticket = *ticket

	return &model.BookingResponse{
		ID:          booking.ID,
		UserID:      booking.UserID,
		EventID:     booking.EventID,
		TicketID:    booking.TicketID,
		Quantity:    booking.Quantity,
		TotalAmount: booking.TotalAmount,
		Status:      booking.Status,
		ExpiredAt:   booking.ExpiredAt,
		CreatedAt:   booking.CreatedAt,
	}, nil
}

func (s *bookingService) GetBookingByID(id uuid.UUID) (*model.BookingResponse, error) {
	booking, err := s.bookingRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("booking not found")
	}

	return &model.BookingResponse{
		ID:          booking.ID,
		UserID:      booking.UserID,
		EventID:     booking.EventID,
		TicketID:    booking.TicketID,
		Quantity:    booking.Quantity,
		TotalAmount: booking.TotalAmount,
		Status:      booking.Status,
		ExpiredAt:   booking.ExpiredAt,
		CreatedAt:   booking.CreatedAt,
	}, nil
}

func (s *bookingService) GetAllBookings() ([]model.BookingResponse, error) {
	bookings, err := s.bookingRepo.FindAll()
	if err != nil {
		return nil, err
	}

	var response []model.BookingResponse
	for _, booking := range bookings {
		response = append(response, model.BookingResponse{
			ID:          booking.ID,
			UserID:      booking.UserID,
			EventID:     booking.EventID,
			TicketID:    booking.TicketID,
			Quantity:    booking.Quantity,
			TotalAmount: booking.TotalAmount,
			Status:      booking.Status,
			ExpiredAt:   booking.ExpiredAt,
			CreatedAt:   booking.CreatedAt,
		})
	}

	return response, nil
}

func (s *bookingService) UpdateBookingStatus(id uuid.UUID, status string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		bookingRepoTx := s.bookingRepo.WithTx(tx)
		ticketRepoTx := s.ticketRepo.WithTx(tx)

		booking, err := bookingRepoTx.FindByIDForUpdate(id)
		if err != nil {
			return errors.New("booking not found")
		}

		if booking.Status != "PENDING" {
			return errors.New("booking is not pending")
		}

		switch status {
		case "CONFIRMED":
			if err := bookingRepoTx.UpdateStatus(id, "CONFIRMED"); err != nil {
				return err
			}
		case "CANCELLED":
			if _, err := ticketRepoTx.FindByIDForUpdate(booking.TicketID); err != nil {
				return errors.New("ticket not found")
			}

			if err := bookingRepoTx.UpdateStatus(id, "CANCELLED"); err != nil {
				return err
			}

			if err := ticketRepoTx.IncreaseQuota(booking.TicketID, booking.Quantity); err != nil {
				return err
			}
		default:
			return errors.New("invalid status")
		}

		return nil
	})
}

package service

import (
	"errors"
	"payment-service/internal/client"
	"payment-service/internal/model"
	"payment-service/internal/repository"
	"time"

	"github.com/google/uuid"
)

type PaymentService interface {
	CreatePayment(bookingID uuid.UUID, amount float64, paymentMethod string) (*model.PaymentResponse, error)
	GetPaymentByID(id uuid.UUID) (*model.PaymentResponse, error)
	GetAllPayments() ([]model.PaymentResponse, error)
	UpdatePaymentStatus(id uuid.UUID, status string) error
	HandlePaymentGatewayWebhook(paymentID uuid.UUID, status string) error
	HandleBookingExpired(bookingID uuid.UUID) error
}

type paymentService struct {
	paymentRepo   repository.PaymentRepository
	bookingClient client.BookingClient
	userClient    client.UserClient
	webhookClient client.WebhookClient
}

func NewPaymentService(
	paymentRepo repository.PaymentRepository,
	bookingClient client.BookingClient,
	userClient client.UserClient,
	webhookClient client.WebhookClient,
) PaymentService {
	return &paymentService{
		paymentRepo:   paymentRepo,
		bookingClient: bookingClient,
		userClient:    userClient,
		webhookClient: webhookClient,
	}
}

func (s *paymentService) CreatePayment(bookingID uuid.UUID, amount float64, paymentMethod string) (*model.PaymentResponse, error) {
	if amount <= 0 {
		return nil, errors.New("amount must be greater than 0")
	}

	validMethods := map[string]bool{"VA": true, "EWALLET": true, "QRIS": true}
	if !validMethods[paymentMethod] {
		return nil, errors.New("invalid payment method. Allowed: VA, EWALLET, QRIS")
	}

	booking, err := s.bookingClient.GetBookingByID(bookingID)
	if err != nil {
		return nil, errors.New("booking not found: " + err.Error())
	}

	existingPayment, _ := s.paymentRepo.FindByBookingID(bookingID)
	if existingPayment != nil {
		return nil, errors.New("payment already exists for this booking")
	}

	paymentExpiry := *booking.ExpiredAt

	payment := &model.Payment{
		BookingID:     bookingID,
		UserID:        booking.UserID,
		Amount:        amount,
		Currency:      "IDR",
		PaymentMethod: paymentMethod,
		Status:        "PENDING",
		ExpiredAt:     &paymentExpiry,
	}

	if err := s.paymentRepo.Create(payment); err != nil {
		return nil, err
	}

	return &model.PaymentResponse{
		ID:            payment.ID,
		BookingID:     payment.BookingID,
		UserID:        payment.UserID,
		Amount:        payment.Amount,
		Currency:      payment.Currency,
		PaymentMethod: payment.PaymentMethod,
		Status:        payment.Status,
		ExpiredAt:     payment.ExpiredAt,
		PaidAt:        payment.PaidAt,
		CreatedAt:     payment.CreatedAt,
		UpdatedAt:     payment.UpdatedAt,
	}, nil
}

func (s *paymentService) GetPaymentByID(id uuid.UUID) (*model.PaymentResponse, error) {
	payment, err := s.paymentRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("payment not found")
	}

	return &model.PaymentResponse{
		ID:            payment.ID,
		BookingID:     payment.BookingID,
		UserID:        payment.UserID,
		Amount:        payment.Amount,
		Currency:      payment.Currency,
		PaymentMethod: payment.PaymentMethod,
		Status:        payment.Status,
		ExpiredAt:     payment.ExpiredAt,
		PaidAt:        payment.PaidAt,
		CreatedAt:     payment.CreatedAt,
		UpdatedAt:     payment.UpdatedAt,
	}, nil
}

func (s *paymentService) GetAllPayments() ([]model.PaymentResponse, error) {
	payments, err := s.paymentRepo.FindAll()
	if err != nil {
		return nil, err
	}

	var response []model.PaymentResponse
	for _, payment := range payments {
		response = append(response, model.PaymentResponse{
			ID:            payment.ID,
			BookingID:     payment.BookingID,
			UserID:        payment.UserID,
			Amount:        payment.Amount,
			Currency:      payment.Currency,
			PaymentMethod: payment.PaymentMethod,
			Status:        payment.Status,
			ExpiredAt:     payment.ExpiredAt,
			PaidAt:        payment.PaidAt,
			CreatedAt:     payment.CreatedAt,
			UpdatedAt:     payment.UpdatedAt,
		})
	}

	return response, nil
}

func (s *paymentService) UpdatePaymentStatus(id uuid.UUID, status string) error {
	payment, err := s.paymentRepo.FindByID(id)
	if err != nil {
		return errors.New("payment not found")
	}

	validStatuses := map[string]bool{"PENDING": true, "PAID": true, "FAILED": true, "EXPIRED": true}
	if !validStatuses[status] {
		return errors.New("invalid status")
	}

	payment.Status = status
	payment.UpdatedAt = time.Now()

	if status == "PAID" && payment.PaidAt == nil {
		now := time.Now()
		payment.PaidAt = &now
	}

	return s.paymentRepo.Update(payment)
}

func (s *paymentService) HandlePaymentGatewayWebhook(paymentID uuid.UUID, status string) error {
	payment, err := s.paymentRepo.FindByID(paymentID)
	if err != nil {
		return errors.New("payment not found")
	}

	if payment.Status != "PENDING" {
		return errors.New("payment is not pending")
	}

	payment.Status = status
	payment.UpdatedAt = time.Now()

	switch status {
	case "PAID":
		now := time.Now()
		payment.PaidAt = &now
	
		if err := s.webhookClient.NotifyBookingService("payment.success", payment.ID, payment.BookingID); err != nil {
				return errors.New("failed to notify booking service: " + err.Error())
		}
	case "FAILED", "EXPIRED":
		if err := s.webhookClient.NotifyBookingService("payment.failed", payment.ID, payment.BookingID); err != nil {
			return errors.New("failed to notify booking service: " + err.Error())
		}
	}

	return s.paymentRepo.Update(payment)
}

func (s *paymentService) HandleBookingExpired(bookingID uuid.UUID) error {
	payment, err := s.paymentRepo.FindByBookingID(bookingID)
	if err != nil {
		// No payment found for this booking, which is okay
		return nil
	}

	// Only expire pending payments
	if payment.Status != "PENDING" {
		return nil
	}

	payment.Status = "EXPIRED"
	payment.UpdatedAt = time.Now()

	return s.paymentRepo.Update(payment)
}

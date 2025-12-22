package handler

import (
	"payment-service/internal/client"
	"payment-service/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type PaymentHandler struct {
	service    service.PaymentService
	userClient client.UserClient
}

func NewPaymentHandler(service service.PaymentService, userClient client.UserClient) *PaymentHandler {
	return &PaymentHandler{
		service:    service,
		userClient: userClient,
	}
}

type CreatePaymentRequest struct {
	BookingID     string  `json:"booking_id" validate:"required"`
	Amount        float64 `json:"amount" validate:"required"`
	PaymentMethod string  `json:"payment_method" validate:"required"`
}

type ProcessPaymentRequest struct {
	PaymentMethod string `json:"payment_method" validate:"required"`
}

type UpdatePaymentStatusRequest struct {
	Status string `json:"status" validate:"required"`
}

type PaymentGatewayWebhookRequest struct {
	PaymentID string `json:"payment_id" validate:"required"`
	Status    string `json:"status" validate:"required"`
}
type BookingWebhookRequest struct {
	Event     string `json:"event"`      // booking.expired, booking.cancelled
	BookingID string `json:"booking_id"`
	Status    string `json:"status"`
}

func (h *PaymentHandler) CreatePayment(c *fiber.Ctx) error {
	var req CreatePaymentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.BookingID == "" || req.Amount <= 0 || req.PaymentMethod == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "booking_id, amount, and payment_method are required and must be valid",
		})
	}

	// Parse booking ID
	bookingID, err := uuid.Parse(req.BookingID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid booking ID format",
		})
	}

	payment, err := h.service.CreatePayment(bookingID, req.Amount, req.PaymentMethod)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Payment created successfully",
		"data":    payment,
	})
}

func (h *PaymentHandler) GetPaymentByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid payment ID",
		})
	}

	payment, err := h.service.GetPaymentByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Payment retrieved successfully",
		"data":    payment,
	})
}

func (h *PaymentHandler) GetAllPayments(c *fiber.Ctx) error {
	payments, err := h.service.GetAllPayments()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Payments retrieved successfully",
		"data":    payments,
	})
}

func (h *PaymentHandler) UpdatePaymentStatus(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid payment ID",
		})
	}

	var req UpdatePaymentStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Status == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "status is required",
		})
	}

	if err := h.service.UpdatePaymentStatus(id, req.Status); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Payment status updated successfully",
	})
}

func (h *PaymentHandler) HandlePaymentGatewayWebhook(c *fiber.Ctx) error {
	var req PaymentGatewayWebhookRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	paymentID, err := uuid.Parse(req.PaymentID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid payment ID format",
		})
	}

	if req.PaymentID == "" || req.Status == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "payment_id and status are required",
		})
	}

	validStatuses := map[string]bool{"PENDING": true, "PAID": true, "FAILED": true, "EXPIRED": true}
	if !validStatuses[req.Status] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid status",
		})
	}

	if err := h.service.HandlePaymentGatewayWebhook(paymentID, req.Status); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Webhook processed successfully",
	})
}

func (h *PaymentHandler) HandleBookingWebhook(c *fiber.Ctx) error {
	var req BookingWebhookRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.BookingID == "" || req.Event == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "booking_id and event are required",
		})
	}

	// Parse booking ID
	bookingID, err := uuid.Parse(req.BookingID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid booking ID format",
		})
	}

	// Handle different booking events
	switch req.Event {
	case "booking.expired", "booking.cancelled":
		if err := h.service.HandleBookingExpired(bookingID); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Unknown event type",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Booking webhook processed successfully",
	})
}

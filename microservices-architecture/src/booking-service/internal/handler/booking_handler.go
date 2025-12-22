package handler

import (
	"booking-service/internal/client"
	"booking-service/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type BookingHandler struct {
	service    service.BookingService
	userClient client.UserClient
}

func NewBookingHandler(service service.BookingService, userClient client.UserClient) *BookingHandler {
	return &BookingHandler{
		service:    service,
		userClient: userClient,
	}
}

type CreateBookingRequest struct {
	EventID  string `json:"event_id" validate:"required"`
	TicketID string `json:"ticket_id" validate:"required"`
	Quantity int    `json:"quantity" validate:"required"`
}

type UpdateBookingStatusRequest struct {
	Status string `json:"status" validate:"required"`
}

type PaymentWebhookRequest struct {
	Event     string `json:"event"`      // payment.success, payment.failed, payment.expired
	PaymentID string `json:"payment_id"`
	BookingID string `json:"booking_id"`
	Status    string `json:"status"`
}

func (h *BookingHandler) CreateBooking(c *fiber.Ctx) error {
	authToken := c.Get("Authorization")
	if authToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization token is required",
		})
	}

	user, err := h.userClient.GetAuthenticatedUser(authToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Failed to authenticate user: " + err.Error(),
		})
	}

	var req CreateBookingRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.EventID == "" || req.TicketID == "" || req.Quantity <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "All fields are required and must be valid",
		})
	}

	eventID, err := uuid.Parse(req.EventID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid event ID format",
		})
	}

	ticketID, err := uuid.Parse(req.TicketID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ticket ID format",
		})
	}

	userID, err := uuid.Parse(user.ID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID format",
		})
	}

	booking, err := h.service.CreateBooking(userID, eventID, ticketID, req.Quantity)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Booking created successfully",
		"data":    booking,
	})
}

func (h *BookingHandler) GetBookingByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid booking ID",
		})
	}

	booking, err := h.service.GetBookingByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Booking retrieved successfully",
		"data":    booking,
	})
}

func (h *BookingHandler) GetAllBookings(c *fiber.Ctx) error {
	bookings, err := h.service.GetAllBookings()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Bookings retrieved successfully",
		"data":    bookings,
	})
}

func (h *BookingHandler) UpdateBookingStatus(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid booking ID",
		})
	}

	var req UpdateBookingStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Status == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Status is required",
		})
	}

	validStatuses := map[string]bool{
		"PENDING":   true,
		"PAID":      true,
		"CONFIRMED": true,
		"CANCELLED": true,
	}

	if !validStatuses[req.Status] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid status",
		})
	}

	if err := h.service.UpdateBookingStatus(id, req.Status); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Booking status updated successfully",
	})
}

func (h *BookingHandler) HandlePaymentWebhook(c *fiber.Ctx) error {
	var req PaymentWebhookRequest
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

	bookingID, err := uuid.Parse(req.BookingID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid booking ID format",
		})
	}

	switch req.Event {
	case "payment.success":
		if err := h.service.UpdateBookingStatus(bookingID, "CONFIRMED"); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	case "payment.failed", "payment.expired":
		if err := h.service.UpdateBookingStatus(bookingID, "CANCELLED"); err != nil {
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
		"message": "Payment webhook processed successfully",
	})
}

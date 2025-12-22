package handler

import (
	"booking-service/internal/model"
	"booking-service/internal/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type TicketHandler struct {
	ticketRepo repository.TicketRepository
}

func NewTicketHandler(ticketRepo repository.TicketRepository) *TicketHandler {
	return &TicketHandler{
		ticketRepo: ticketRepo,
	}
}

func (h *TicketHandler) GetTicketByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ticket ID",
		})
	}

	ticket, err := h.ticketRepo.FindByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Ticket not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Ticket retrieved successfully",
		"data": model.TicketResponse{
			ID:       ticket.ID,
			EventID:  ticket.EventID,
			Category: ticket.Category,
			Price:    ticket.Price,
			Quota:    ticket.Quota,
		},
	})
}

func (h *TicketHandler) GetAllTickets(c *fiber.Ctx) error {
	tickets, err := h.ticketRepo.FindAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve tickets",
		})
	}

	var response []model.TicketResponse
	for _, ticket := range tickets {
		response = append(response, model.TicketResponse{
			ID:       ticket.ID,
			EventID:  ticket.EventID,
			Category: ticket.Category,
			Price:    ticket.Price,
			Quota:    ticket.Quota,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Tickets retrieved successfully",
		"data":    response,
	})
}

func (h *TicketHandler) GetTicketsByEventID(c *fiber.Ctx) error {
	eventIDParam := c.Params("id")
	eventID, err := uuid.Parse(eventIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid event ID",
		})
	}

	tickets, err := h.ticketRepo.FindByEventID(eventID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve tickets",
		})
	}

	var response []model.TicketResponse
	for _, ticket := range tickets {
		response = append(response, model.TicketResponse{
			ID:       ticket.ID,
			EventID:  ticket.EventID,
			Category: ticket.Category,
			Price:    ticket.Price,
			Quota:    ticket.Quota,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Tickets retrieved successfully",
		"data":    response,
	})
}

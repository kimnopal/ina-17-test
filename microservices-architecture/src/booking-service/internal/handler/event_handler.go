package handler

import (
	"booking-service/internal/model"
	"booking-service/internal/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type EventHandler struct {
	eventRepo repository.EventRepository
}

func NewEventHandler(eventRepo repository.EventRepository) *EventHandler {
	return &EventHandler{
		eventRepo: eventRepo,
	}
}

func (h *EventHandler) GetEventByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid event ID",
		})
	}

	event, err := h.eventRepo.FindByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Event not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Event retrieved successfully",
		"data": model.EventResponse{
			ID:          event.ID,
			Name:        event.Name,
			Description: event.Description,
			EventDate:   event.EventDate,
			CreatedAt:   event.CreatedAt,
		},
	})
}

func (h *EventHandler) GetAllEvents(c *fiber.Ctx) error {
	events, err := h.eventRepo.FindAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve events",
		})
	}

	var response []model.EventResponse
	for _, event := range events {
		response = append(response, model.EventResponse{
			ID:          event.ID,
			Name:        event.Name,
			Description: event.Description,
			EventDate:   event.EventDate,
			CreatedAt:   event.CreatedAt,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Events retrieved successfully",
		"data":    response,
	})
}

package migrations

import (
	"booking-service/internal/model"
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func SeedData(db *gorm.DB) {
	// Check if data already exists
	var eventCount int64
	db.Model(&model.Event{}).Count(&eventCount)
	if eventCount > 0 {
		log.Println("Data already seeded, skipping...")
		return
	}

	// Seed Events
	events := []model.Event{
		{
			ID:          uuid.New(),
			Name:        "Rock Concert 2025",
			Description: "An amazing rock concert featuring top bands from around the world",
			EventDate:   time.Date(2025, 3, 15, 19, 0, 0, 0, time.UTC),
		},
		{
			ID:          uuid.New(),
			Name:        "Tech Conference 2025",
			Description: "Annual technology conference with keynote speakers and workshops",
			EventDate:   time.Date(2025, 4, 20, 9, 0, 0, 0, time.UTC),
		},
		{
			ID:          uuid.New(),
			Name:        "Jazz Festival",
			Description: "Three-day jazz festival featuring international and local artists",
			EventDate:   time.Date(2025, 5, 10, 18, 0, 0, 0, time.UTC),
		},
		{
			ID:          uuid.New(),
			Name:        "Food & Wine Expo",
			Description: "Culinary experience with renowned chefs and wine tasting",
			EventDate:   time.Date(2025, 6, 5, 12, 0, 0, 0, time.UTC),
		},
		{
			ID:          uuid.New(),
			Name:        "Summer Music Festival",
			Description: "Outdoor music festival with multiple stages and diverse genres",
			EventDate:   time.Date(2025, 7, 25, 15, 0, 0, 0, time.UTC),
		},
	}

	for _, event := range events {
		if err := db.Create(&event).Error; err != nil {
			log.Printf("Failed to seed event %s: %v", event.Name, err)
		}
	}
	log.Println("Events seeded successfully")

	// Seed Tickets for each event
	tickets := []model.Ticket{
		// Rock Concert 2025 tickets
		{
			ID:       uuid.New(),
			EventID:  events[0].ID,
			Category: "VIP",
			Price:    150.00,
			Quota:    100,
		},
		{
			ID:       uuid.New(),
			EventID:  events[0].ID,
			Category: "Regular",
			Price:    75.00,
			Quota:    500,
		},
		// Tech Conference 2025 tickets
		{
			ID:       uuid.New(),
			EventID:  events[1].ID,
			Category: "VIP",
			Price:    200.00,
			Quota:    50,
		},
		{
			ID:       uuid.New(),
			EventID:  events[1].ID,
			Category: "Regular",
			Price:    100.00,
			Quota:    300,
		},
		// Jazz Festival tickets
		{
			ID:       uuid.New(),
			EventID:  events[2].ID,
			Category: "VIP",
			Price:    180.00,
			Quota:    80,
		},
		{
			ID:       uuid.New(),
			EventID:  events[2].ID,
			Category: "Regular",
			Price:    90.00,
			Quota:    400,
		},
		// Food & Wine Expo tickets
		{
			ID:       uuid.New(),
			EventID:  events[3].ID,
			Category: "VIP",
			Price:    120.00,
			Quota:    60,
		},
		{
			ID:       uuid.New(),
			EventID:  events[3].ID,
			Category: "Regular",
			Price:    60.00,
			Quota:    250,
		},
		// Summer Music Festival tickets
		{
			ID:       uuid.New(),
			EventID:  events[4].ID,
			Category: "VIP",
			Price:    250.00,
			Quota:    150,
		},
		{
			ID:       uuid.New(),
			EventID:  events[4].ID,
			Category: "Regular",
			Price:    100.00,
			Quota:    1000,
		},
	}

	for _, ticket := range tickets {
		if err := db.Create(&ticket).Error; err != nil {
			log.Printf("Failed to seed ticket for event %s: %v", ticket.EventID, err)
		}
	}
	log.Println("Tickets seeded successfully")
}

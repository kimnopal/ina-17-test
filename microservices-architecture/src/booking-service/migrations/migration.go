package migrations

import (
	"booking-service/internal/model"
	"log"

	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) {
	err := db.AutoMigrate(
		&model.Event{},
		&model.Ticket{},
		&model.Booking{},
	)
	if err != nil {
		log.Fatal("Failed to run migrations:", err)
	}
	log.Println("Migrations completed successfully")

	// Run seeders
	SeedData(db)
}

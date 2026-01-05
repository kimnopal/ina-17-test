package main

import (
	"booking-service/config"
	"booking-service/internal/client"
	"booking-service/internal/handler"
	"booking-service/internal/repository"
	"booking-service/internal/service"
	"booking-service/migrations"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// Connect to database
	config.ConnectDatabase()

	// Run migrations
	migrations.RunMigrations(config.DB)

	// Initialize clients
	userClient := client.NewUserClient()
	paymentClient := client.NewPaymentClient()
	webhookClient := client.NewWebhookClient()

	// Initialize repositories
	bookingRepo := repository.NewBookingRepository(config.DB)
	eventRepo := repository.NewEventRepository(config.DB)
	ticketRepo := repository.NewTicketRepository(config.DB)

	// Initialize services
	bookingService := service.NewBookingService(config.DB, bookingRepo, ticketRepo, eventRepo, userClient, paymentClient, webhookClient)

	// Initialize handlers
	bookingHandler := handler.NewBookingHandler(bookingService, userClient)
	eventHandler := handler.NewEventHandler(eventRepo)
	ticketHandler := handler.NewTicketHandler(ticketRepo)

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName: "Booking Service",
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New())

	// Routes
	api := app.Group("/api/v1")

	// Event routes
	events := api.Group("/events")
	events.Get("/", eventHandler.GetAllEvents)
	events.Get("/:id", eventHandler.GetEventByID)
	events.Get("/:id/tickets", ticketHandler.GetTicketsByEventID)

	// Ticket routes
	tickets := api.Group("/tickets")
	tickets.Get("/:id", ticketHandler.GetTicketByID)

	// Booking routes
	bookings := api.Group("/bookings")
	bookings.Post("/", bookingHandler.CreateBooking)
	bookings.Get("/", bookingHandler.GetAllBookings)
	bookings.Get("/:id", bookingHandler.GetBookingByID)
	bookings.Put("/:id/status", bookingHandler.UpdateBookingStatus)
	bookings.Post("/webhook/payment", bookingHandler.HandlePaymentWebhook)

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "booking-service",
		})
	})

	// Start server
	port := ":3002"
	log.Printf("Server starting on port %s", port)
	if err := app.Listen(port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

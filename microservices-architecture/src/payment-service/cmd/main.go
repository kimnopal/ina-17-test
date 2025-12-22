package main

import (
	"log"
	"payment-service/config"
	"payment-service/internal/client"
	"payment-service/internal/handler"
	"payment-service/internal/repository"
	"payment-service/internal/service"
	"payment-service/migrations"

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
	bookingClient := client.NewBookingClient()
	webhookClient := client.NewWebhookClient()

	// Initialize repositories
	paymentRepo := repository.NewPaymentRepository(config.DB)

	// Initialize services
	paymentService := service.NewPaymentService(paymentRepo, bookingClient, userClient, webhookClient)

	// Initialize handlers
	paymentHandler := handler.NewPaymentHandler(paymentService, userClient)

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName: "Payment Service",
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New())

	// Routes
	api := app.Group("/api/v1")

	// Payment routes
	payments := api.Group("/payments")
	payments.Post("/webhook/payment-gateway", paymentHandler.HandlePaymentGatewayWebhook)
	// payments.Post("/webhook/booking", paymentHandler.HandleBookingWebhook) // Webhook from booking service
	payments.Post("/", paymentHandler.CreatePayment)
	payments.Get("/", paymentHandler.GetAllPayments)
	payments.Get("/:id", paymentHandler.GetPaymentByID)
	payments.Put("/:id/status", paymentHandler.UpdatePaymentStatus)

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "payment-service",
		})
	})

	port := ":3002"
	log.Printf("Server starting on port %s", port)
	if err := app.Listen(port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

package main

import (
	"log"
	"user-service/config"
	"user-service/internal/handler"
	"user-service/internal/repository"
	"user-service/internal/service"
	"user-service/migrations"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// Connect to database
	config.ConnectDatabase()

	// Run migrations
	migrations.RunMigrations(config.DB)

	// Initialize repository, service, and handler
	userRepo := repository.NewUserRepository(config.DB)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName: "User Service",
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New())

	// Routes
	api := app.Group("/api/v1")

	// User routes
	api.Post("/login", userHandler.Login)
	users := api.Group("/users")
	users.Post("/", userHandler.CreateUser)
	users.Get("/auth", userHandler.GetAuthenticatedUser)

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "user-service",
		})
	})

	// Start server
	port := ":3000"
	log.Printf("Server starting on port %s", port)
	if err := app.Listen(port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

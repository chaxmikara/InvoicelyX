package routes

import (
	"InvoicelyX/handlers"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRoutes(app *fiber.App, db *mongo.Database) {
	// Initialize handlers
	userHandler := handlers.NewUserHandler(db)

	// API routes
	api := app.Group("/api")

	// User routes
	api.Post("/users", userHandler.CreateUser)
}

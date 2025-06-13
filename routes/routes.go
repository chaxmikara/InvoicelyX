package routes

import (
	"InvoicelyX/handlers"
	middleware "InvoicelyX/middlewares"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRoutes(app *fiber.App, db *mongo.Database) {
	userHandler := handlers.NewUserHandler(db)

	api := app.Group("/api")

	// User routes
	api.Post("/register", userHandler.CreateUser)
	api.Post("/users/login", userHandler.Login)
	// Protected routes (authentication required)
	protected := api.Group("/", middleware.JWTMiddleware())
	//protected.Post("/users/logout", userHandler.Logout)
	protected.Get("/users/profile", userHandler.GetProfile)
}

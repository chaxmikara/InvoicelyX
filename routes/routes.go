package routes

import (
	"InvoicelyX/handlers"
	middleware "InvoicelyX/middlewares"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRoutes(app *fiber.App, db *mongo.Database) {
	// Initialize handlers
	userHandler := handlers.NewUserHandler(db)
	invoiceHandler := handlers.NewInvoiceHandler(db)

	api := app.Group("/api")

	// Public routes
	api.Post("/users", userHandler.CreateUser)  // Register user
	api.Post("/users/login", userHandler.Login) // Login user

	// Protected routes (authentication required)
	protected := api.Group("/", middleware.JWTMiddleware())
	protected.Post("/users/logout", userHandler.Logout)     // Logout user
	protected.Get("/users/profile", userHandler.GetProfile) // Get profile

	// Invoice routes (protected)
	protected.Post("/invoices", invoiceHandler.CreateInvoice)       // Create invoice
	protected.Get("/invoices", invoiceHandler.GetAllInvoices)       // Get all invoices
	protected.Get("/invoices/:id", invoiceHandler.GetInvoiceByID)   // Get invoice by ID
	protected.Put("/invoices/:id", invoiceHandler.UpdateInvoice)    // Update invoice
	protected.Delete("/invoices/:id", invoiceHandler.DeleteInvoice) // Delete invoice

	// Admin only routes

}

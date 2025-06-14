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
	protected.Post("/users/logout", userHandler.Logout)                 // Logout user
	protected.Get("/users/profile", userHandler.GetProfile)             // Get profile
	protected.Put("/users/profile", userHandler.UpdateProfile)          // Update profile
	protected.Put("/users/change-password", userHandler.ChangePassword) // Change password

	// User Invoice routes (user can only access their own invoices)
	protected.Post("/invoices", invoiceHandler.CreateInvoice)             // Create invoice
	protected.Get("/my-invoices", invoiceHandler.GetMyInvoices)           // Get user's invoices with filtering
	protected.Get("/invoices/:id", invoiceHandler.GetInvoiceByID)         // Get user's invoice by ID
	protected.Put("/invoices/:id", invoiceHandler.UpdateInvoice)          // Update user's invoice
	protected.Delete("/invoices/:id", invoiceHandler.DeleteInvoice)       // Delete user's invoice
	protected.Get("/invoices/:id/pdf", invoiceHandler.GenerateInvoicePDF) // Generate PDF
	protected.Get("/my-invoices/export", invoiceHandler.ExportMyInvoices) // Export user's invoices to CSV

	// Admin Invoice routes (admin can access all invoices)
	adminInvoices := protected.Group("/admin/invoices", middleware.RequireRole("admin"))
	adminInvoices.Get("/", invoiceHandler.GetAllInvoices) // Admin: Get all invoices

	// Admin only routes
	adminOnly := protected.Group("/", middleware.RequireRole("admin"))

	// User Management
	adminOnly.Get("/admin/users", userHandler.GetAllUsers)                   // Get all users
	adminOnly.Put("/admin/users/:id/activate", userHandler.ActivateUser)     // Activate user
	adminOnly.Put("/admin/users/:id/deactivate", userHandler.DeactivateUser) // Deactivate user
	adminOnly.Delete("/admin/users/:id", userHandler.DeleteUser)             // Delete user

	// System Statistics & Monitoring
	adminOnly.Get("/admin/stats/invoices", invoiceHandler.GetInvoiceStatistics) // Invoice statistics
	adminOnly.Get("/admin/system/stats", userHandler.GetSystemStats)            // System statistics
	adminOnly.Get("/admin/system/health", userHandler.GetSystemHealth)          // System health check

	// Data Export
	adminOnly.Get("/admin/export/invoices", invoiceHandler.ExportAllInvoices) // Export all invoices
	adminOnly.Get("/admin/export/users", userHandler.ExportUserData)          // Export user data
}

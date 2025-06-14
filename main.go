package main

import (
	"InvoicelyX/config"
	"InvoicelyX/routes"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect to MongoDB
	client := config.ConnectDB()
	defer func() {
		err := client.Disconnect(context.TODO())
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Disconnected from MongoDB.")
	}()

	// Get database
	db := client.Database("invoicelyx")

	// Create indexes
	err = config.CreateIndexes(db)
	if err != nil {
		log.Printf("Warning: Failed to create indexes: %v", err)
	}

	// Create Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		},
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New())

	// Setup routes
	routes.SetupRoutes(app, db)

	// Health check route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "InvoicelyX API is running!",
			"status":  "healthy",
		})
	})

	// Get port from environment variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	fmt.Printf("Server starting on port %s\n", port)

	// Start the server
	log.Fatal(app.Listen(":" + port))
}

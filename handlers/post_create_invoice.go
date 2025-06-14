package handlers

import (
	"InvoicelyX/models"
	"InvoicelyX/utils"
	"context"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// CreateInvoice creates a new invoice
func (h *InvoiceHandler) CreateInvoice(c *fiber.Ctx) error {
	var invoice models.Invoice

	if err := c.BodyParser(&invoice); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	// Validate the invoice
	validate := validator.New()
	if err := validate.Struct(&invoice); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Validation error",
			"errors":  err.Error(),
		})
	}

	// Additional validation for items
	if !utils.ValidateInvoiceItems(invoice.Items) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid invoice items",
		})
	}
	// Generate UUID for invoice (this is primary ID)
	invoice.InvoiceID = uuid.New().String()

	// Associate invoice with current user
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Invalid user session",
		})
	}
	invoice.UserID = userID

	// Calculate total if not provided
	if invoice.Total == 0 {
		invoice.Total = utils.CalculateInvoiceTotal(invoice.Items)
	}

	// Set timestamps
	invoice.CreatedAt = time.Now()
	invoice.UpdatedAt = time.Now()

	// Create invoice in database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := h.DB.Collection("invoices")

	_, err := collection.InsertOne(ctx, invoice)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to create invoice",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Invoice created successfully",
		"data":    invoice,
	})
}

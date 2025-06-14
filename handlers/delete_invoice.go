package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

// DeleteInvoice deletes an invoice by ID
func (h *InvoiceHandler) DeleteInvoice(c *fiber.Ctx) error {
	invoiceID := c.Params("id")
	if invoiceID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invoice ID is required",
		})
	}

	// Get user ID from JWT
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Invalid user session",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := h.DB.Collection("invoices")

	// Delete only user's own invoice
	result, err := collection.DeleteOne(ctx, bson.M{
		"invoice_id": invoiceID,
		"user_id":    userID,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to delete invoice",
			"error":   err.Error(),
		})
	}

	if result.DeletedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Invoice not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Invoice deleted successfully",
	})
}

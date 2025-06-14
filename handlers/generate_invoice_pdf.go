package handlers

import (
	"InvoicelyX/models"
	"InvoicelyX/utils"
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GenerateInvoicePDF generates and returns a PDF for a specific invoice
func (h *InvoiceHandler) GenerateInvoicePDF(c *fiber.Ctx) error {
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

	// Find invoice by ID and user ID (ensure user owns the invoice)
	var invoice models.Invoice
	err := collection.FindOne(ctx, bson.M{
		"invoice_id": invoiceID,
		"user_id":    userID,
	}).Decode(&invoice)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Invoice not found or access denied",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Database error",
			"error":   err.Error(),
		})
	}

	// Generate PDF
	pdfBytes, err := utils.GenerateInvoicePDF(&invoice)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to generate PDF",
			"error":   err.Error(),
		})
	}

	// Set headers for PDF download
	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", "attachment; filename=invoice_"+invoiceID+".pdf")
	c.Set("Content-Length", string(rune(len(pdfBytes))))

	return c.Send(pdfBytes)
}

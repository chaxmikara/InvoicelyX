package handlers

import (
	"InvoicelyX/models"
	"InvoicelyX/utils"
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

// ExportMyInvoices exports user's invoices to CSV
func (h *InvoiceHandler) ExportMyInvoices(c *fiber.Ctx) error {
	// Get user ID from JWT
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Invalid user session",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	collection := h.DB.Collection("invoices")

	// Find all invoices for this user
	cursor, err := collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to fetch invoices",
			"error":   err.Error(),
		})
	}
	defer cursor.Close(ctx)

	var invoices []models.Invoice
	if err = cursor.All(ctx, &invoices); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to decode invoices",
			"error":   err.Error(),
		})
	}

	// Generate CSV
	csvData, err := utils.GenerateInvoicesCSV(invoices)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to generate CSV",
			"error":   err.Error(),
		})
	}

	// Set headers for CSV download
	c.Set("Content-Type", "text/csv")
	c.Set("Content-Disposition", "attachment; filename=my_invoices_"+time.Now().Format("2006-01-02")+".csv")

	return c.SendString(string(csvData))
}

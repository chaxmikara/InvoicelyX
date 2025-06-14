package handlers

import (
	"InvoicelyX/models"
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

// GetAllInvoices returns all invoices
func (h *InvoiceHandler) GetAllInvoices(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := h.DB.Collection("invoices")

	// Find all invoices
	cursor, err := collection.Find(ctx, bson.M{})
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

	// If no invoices found, return empty array instead of null
	if invoices == nil {
		invoices = []models.Invoice{}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Invoices retrieved successfully",
		"data": fiber.Map{
			"invoices": invoices,
			"count":    len(invoices),
		},
	})
}

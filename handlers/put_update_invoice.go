package handlers

import (
	"InvoicelyX/models"
	"InvoicelyX/utils"
	"context"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

// UpdateInvoice updates an existing invoice
func (h *InvoiceHandler) UpdateInvoice(c *fiber.Ctx) error {
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

	var updateData models.Invoice
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}
	// Validate the update data
	validate := validator.New()
	if err := validate.Struct(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Validation error",
			"errors":  err.Error(),
		})
	}

	// Additional validation for items
	if len(updateData.Items) > 0 && !utils.ValidateInvoiceItems(updateData.Items) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid invoice items",
		})
	}

	// Recalculate total if items are provided
	if len(updateData.Items) > 0 {
		updateData.Total = utils.CalculateInvoiceTotal(updateData.Items)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := h.DB.Collection("invoices")
	// Create update document
	updateDoc := bson.M{
		"$set": bson.M{
			"customer":   updateData.Customer,
			"email":      updateData.Email,
			"items":      updateData.Items,
			"total":      updateData.Total,
			"updated_at": time.Now(),
		},
	}

	// Update only user's own invoice
	result, err := collection.UpdateOne(ctx, bson.M{
		"invoice_id": invoiceID,
		"user_id":    userID,
	}, updateDoc)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to update invoice",
			"error":   err.Error(),
		})
	}

	if result.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Invoice not found",
		})
	}
	// Get updated invoice
	var updatedInvoice models.Invoice
	err = collection.FindOne(ctx, bson.M{
		"invoice_id": invoiceID,
		"user_id":    userID,
	}).Decode(&updatedInvoice)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to retrieve updated invoice",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Invoice updated successfully",
		"data":    updatedInvoice,
	})
}

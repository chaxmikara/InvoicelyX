package handlers

import (
	"InvoicelyX/models"
	"context"
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

// ExportAllInvoices - Admin endpoint to export all invoices for compliance
func (h *InvoiceHandler) ExportAllInvoices(c *fiber.Ctx) error {
	// Get export format from query parameter (default: csv)
	format := c.Query("format", "csv")
	if format != "csv" && format != "json" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid format. Supported formats: csv, json",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	collection := h.DB.Collection("invoices")

	// Get all invoices
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

	if format == "json" {
		return h.exportInvoicesAsJSON(c, invoices)
	}

	return h.exportInvoicesAsCSV(c, invoices)
}

// exportInvoicesAsJSON exports invoices in JSON format
func (h *InvoiceHandler) exportInvoicesAsJSON(c *fiber.Ctx, invoices []models.Invoice) error {
	filename := fmt.Sprintf("invoices_export_%s.json", time.Now().Format("2006-01-02_15-04-05"))

	c.Set("Content-Type", "application/json")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	exportData := fiber.Map{
		"export_metadata": fiber.Map{
			"exported_at":    time.Now(),
			"total_invoices": len(invoices),
			"exported_by":    c.Locals("email"), // From JWT middleware
			"format":         "json",
		},
		"invoices": invoices,
	}

	return c.JSON(exportData)
}

// exportInvoicesAsCSV exports invoices in CSV format
func (h *InvoiceHandler) exportInvoicesAsCSV(c *fiber.Ctx, invoices []models.Invoice) error {
	filename := fmt.Sprintf("invoices_export_%s.csv", time.Now().Format("2006-01-02_15-04-05"))

	c.Set("Content-Type", "text/csv")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	// Create CSV content
	var csvContent strings.Builder
	writer := csv.NewWriter(&csvContent)

	// Write header
	header := []string{
		"Invoice ID", "Customer", "Email", "Total", "Items Count",
		"Created At", "Updated At", "Item Details",
	}
	writer.Write(header)

	// Write data
	for _, invoice := range invoices {
		// Format item details
		var itemDetails []string
		for _, item := range invoice.Items {
			itemDetail := fmt.Sprintf("%s (Qty: %d, Price: %.2f)",
				item.Name, item.Quantity, item.Price)
			itemDetails = append(itemDetails, itemDetail)
		}

		record := []string{
			invoice.InvoiceID,
			invoice.Customer,
			invoice.Email,
			strconv.FormatFloat(invoice.Total, 'f', 2, 64),
			strconv.Itoa(len(invoice.Items)),
			invoice.CreatedAt.Format("2006-01-02 15:04:05"),
			invoice.UpdatedAt.Format("2006-01-02 15:04:05"),
			strings.Join(itemDetails, "; "),
		}
		writer.Write(record)
	}

	writer.Flush()

	if err := writer.Error(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to generate CSV",
			"error":   err.Error(),
		})
	}

	return c.Send([]byte(csvContent.String()))
}

// ExportUserData - Admin endpoint to export user data
func (h *UserHandler) ExportUserData(c *fiber.Ctx) error {
	format := c.Query("format", "csv")
	if format != "csv" && format != "json" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid format. Supported formats: csv, json",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	collection := h.DB.Collection("users")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to fetch users",
			"error":   err.Error(),
		})
	}
	defer cursor.Close(ctx)

	var users []models.User
	if err = cursor.All(ctx, &users); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to decode users",
			"error":   err.Error(),
		})
	}

	if format == "json" {
		filename := fmt.Sprintf("users_export_%s.json", time.Now().Format("2006-01-02_15-04-05"))
		c.Set("Content-Type", "application/json")
		c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

		exportData := fiber.Map{
			"export_metadata": fiber.Map{
				"exported_at": time.Now(),
				"total_users": len(users),
				"exported_by": c.Locals("email"),
				"format":      "json",
			},
			"users": users,
		}

		return c.JSON(exportData)
	}

	// CSV export
	filename := fmt.Sprintf("users_export_%s.csv", time.Now().Format("2006-01-02_15-04-05"))
	c.Set("Content-Type", "text/csv")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	var csvContent strings.Builder
	writer := csv.NewWriter(&csvContent)

	// Header
	header := []string{"User ID", "First Name", "Last Name", "Email", "Role", "Is Active", "Created At", "Updated At"}
	writer.Write(header)

	// Data
	for _, user := range users {
		record := []string{
			user.UserID,
			user.FirstName,
			user.LastName,
			user.Email,
			user.Role,
			strconv.FormatBool(user.IsActive),
			user.CreatedAt.Format("2006-01-02 15:04:05"),
			user.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
		writer.Write(record)
	}

	writer.Flush()
	return c.Send([]byte(csvContent.String()))
}

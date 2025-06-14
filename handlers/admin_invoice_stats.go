package handlers

import (
	"InvoicelyX/models"
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

// GetInvoiceStatistics - Admin endpoint for system-wide invoice statistics
func (h *InvoiceHandler) GetInvoiceStatistics(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := h.DB.Collection("invoices")

	// Get all invoices for statistics
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

	// Calculate statistics
	totalInvoices := len(invoices)
	var totalRevenue float64
	var monthlyStats = make(map[string]float64)
	var dailyStats = make(map[string]int)

	now := time.Now()
	thisMonth := now.Format("2006-01")
	today := now.Format("2006-01-02")

	for _, invoice := range invoices {
		totalRevenue += invoice.Total

		// Monthly statistics
		invoiceMonth := invoice.CreatedAt.Format("2006-01")
		monthlyStats[invoiceMonth] += invoice.Total

		// Daily statistics (last 30 days)
		invoiceDay := invoice.CreatedAt.Format("2006-01-02")
		if invoice.CreatedAt.After(now.AddDate(0, 0, -30)) {
			dailyStats[invoiceDay]++
		}
	}

	// Current month revenue
	currentMonthRevenue := monthlyStats[thisMonth]

	// Today's invoice count
	todayCount := dailyStats[today]

	// Average invoice value
	var avgInvoiceValue float64
	if totalInvoices > 0 {
		avgInvoiceValue = totalRevenue / float64(totalInvoices)
	}

	// Get user count for additional stats
	userCollection := h.DB.Collection("users")
	userCount, err := userCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		userCount = 0 // Continue even if user count fails
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Invoice statistics retrieved successfully",
		"data": fiber.Map{
			"overview": fiber.Map{
				"total_invoices":        totalInvoices,
				"total_revenue":         totalRevenue,
				"current_month_revenue": currentMonthRevenue,
				"today_invoices":        todayCount,
				"average_invoice_value": avgInvoiceValue,
				"total_users":           userCount,
			},
			"monthly_revenue": monthlyStats,
			"daily_counts":    dailyStats,
			"generated_at":    time.Now(),
		},
	})
}

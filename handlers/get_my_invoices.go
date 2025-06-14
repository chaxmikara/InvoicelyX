package handlers

import (
	"InvoicelyX/models"
	"context"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

// GetMyInvoices returns invoices for the authenticated user with filtering
func (h *InvoiceHandler) GetMyInvoices(c *fiber.Ctx) error {
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

	// Build query filter
	filter := bson.M{"user_id": userID}

	// Add optional filters from query parameters
	if customer := c.Query("customer"); customer != "" {
		filter["customer"] = bson.M{"$regex": customer, "$options": "i"}
	}

	if email := c.Query("email"); email != "" {
		filter["email"] = bson.M{"$regex": email, "$options": "i"}
	}

	// Date range filtering
	if startDate := c.Query("start_date"); startDate != "" {
		if endDate := c.Query("end_date"); endDate != "" {
			// Parse dates (expecting YYYY-MM-DD format)
			start, err1 := time.Parse("2006-01-02", startDate)
			end, err2 := time.Parse("2006-01-02", endDate)

			if err1 == nil && err2 == nil {
				// Set end date to end of day
				end = end.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
				filter["created_at"] = bson.M{
					"$gte": start,
					"$lte": end,
				}
			}
		}
	}

	// Pagination
	page := 1
	limit := 10

	if p := c.Query("page"); p != "" {
		if pageNum, err := strconv.Atoi(p); err == nil && pageNum > 0 {
			page = pageNum
		}
	}

	if l := c.Query("limit"); l != "" {
		if limitNum, err := strconv.Atoi(l); err == nil && limitNum > 0 && limitNum <= 100 {
			limit = limitNum
		}
	}

	skip := (page - 1) * limit

	// Find invoices with pagination
	cursor, err := collection.Find(ctx, filter, nil)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to fetch invoices",
			"error":   err.Error(),
		})
	}
	defer cursor.Close(ctx)

	var allInvoices []models.Invoice
	if err = cursor.All(ctx, &allInvoices); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to decode invoices",
			"error":   err.Error(),
		})
	}

	// Apply pagination in memory (for simplicity)
	totalCount := len(allInvoices)
	start := skip
	end := skip + limit

	if start > totalCount {
		start = totalCount
	}
	if end > totalCount {
		end = totalCount
	}

	paginatedInvoices := allInvoices[start:end]

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Invoices retrieved successfully",
		"data": fiber.Map{
			"invoices": paginatedInvoices,
			"pagination": fiber.Map{
				"current_page": page,
				"per_page":     limit,
				"total":        totalCount,
				"total_pages":  (totalCount + limit - 1) / limit,
			},
		},
	})
}

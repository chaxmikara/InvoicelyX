package handlers

import (
	"InvoicelyX/models"
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

// GetAllUsers - Admin endpoint to get all users
func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := h.DB.Collection("users")

	// Find all users
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

	// If no users found, return empty array
	if users == nil {
		users = []models.User{}
	}

	// Count active/inactive users
	activeCount := 0
	inactiveCount := 0
	for _, user := range users {
		if user.IsActive {
			activeCount++
		} else {
			inactiveCount++
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Users retrieved successfully",
		"data": fiber.Map{
			"users":          users,
			"total_count":    len(users),
			"active_count":   activeCount,
			"inactive_count": inactiveCount,
		},
	})
}

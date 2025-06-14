package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	// Get user info from JWT (with safety checks)
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Invalid user session",
		})
	}

	email, _ := c.Locals("email").(string)
	role, _ := c.Locals("role").(string)
	firstName, _ := c.Locals("first_name").(string)
	lastName, _ := c.Locals("last_name").(string)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Profile retrieved successfully",
		"data": map[string]interface{}{
			"user_id":   userID,
			"email":     email,
			"role":      role,
			"firstName": firstName,
			"lastName":  lastName,
		},
	})
}

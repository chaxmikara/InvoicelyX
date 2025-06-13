package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	// Get user info from JWT
	userID := c.Locals("user_id").(string)
	email := c.Locals("email").(string)
	role := c.Locals("role").(string)
	firstName := c.Locals("first_name").(string)
	lastName := c.Locals("last_name").(string)

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

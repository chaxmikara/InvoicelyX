package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func (h *UserHandler) Logout(c *fiber.Ctx) error {

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Logout successful. Please remove the token from client storage.",
	})
}

package handlers

import (
	"InvoicelyX/models"
	"InvoicelyX/utils"
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// Login function with JWT
func (h *UserHandler) Login(c *fiber.Ctx) error {
	// Parse request body
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	// Validate request
	if err := h.Validate.Struct(&req); err != nil {
		return h.sendValidationErrorResponse(c, err)
	}

	// Find user by email
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := h.DB.Collection("users")

	var user models.User
	err := collection.FindOne(ctx, bson.M{"email": req.Email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Invalid email or password",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Database error",
			"error":   err.Error(),
		})
	}

	// Check if user is active
	if !user.IsActive {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"message": "Account is deactivated. Please contact administrator.",
		})
	}

	// Compare password with hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Invalid email or password",
		})
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.UserID, user.Email, user.Role, user.FirstName, user.LastName)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to generate token",
			"error":   err.Error(),
		})
	}

	// Update last login time
	updateFilter := bson.M{"_id": user.ID}
	updateData := bson.M{
		"$set": bson.M{
			"lastLoginAt": time.Now(),
			"updatedAt":   time.Now(),
		},
	}

	_, err = collection.UpdateOne(ctx, updateFilter, updateData)
	if err != nil {
		// Log error ...
	}

	loginResponse := map[string]interface{}{
		"user": map[string]interface{}{
			"id":        user.ID,
			"user_id":   user.UserID,
			"firstName": user.FirstName,
			"lastName":  user.LastName,
			"email":     user.Email,
			"role":      user.Role,
			"isActive":  user.IsActive,
			"createdAt": user.CreatedAt,
		},
		"token": token,
	}

	// Send success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Login successful",
		"data":    loginResponse,
	})
}

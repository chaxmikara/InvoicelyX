package handlers

import (
	"InvoicelyX/models"
	"context"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	DB       *mongo.Database
	Validate *validator.Validate
}

func NewUserHandler(db *mongo.Database) *UserHandler {
	return &UserHandler{
		DB:       db,
		Validate: validator.New(),
	}
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	// Use the same User struct for request
	var req models.User
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

	// Additional password validation
	if err := h.validatePassword(req.Password); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := h.DB.Collection("users")

	var existingUser models.User
	err := collection.FindOne(ctx, bson.M{"email": req.Email}).Decode(&existingUser)
	if err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"success": false,
			"message": "User with this email already exists",
		})
	} else if err != mongo.ErrNoDocuments {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Database error",
			"error":   err.Error(),
		})
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to hash password",
			"error":   err.Error(),
		})
	}

	// Generate UUID and set defaults
	req.UserID = uuid.New().String()
	req.Password = string(hashedPassword)
	if req.Role == "" {
		req.Role = "user"
	}
	req.IsActive = true
	req.CreatedAt = time.Now()
	req.UpdatedAt = time.Now()

	// Insert user
	result, err := collection.InsertOne(ctx, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to create user",
			"error":   err.Error(),
		})
	}

	// Clear password before response
	req.Password = ""

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "User created successfully",
		"data":    req,
		"id":      result.InsertedID,
	})
}

func (h *UserHandler) validatePassword(password string) error {
	if len(password) < 8 {
		return fiber.NewError(fiber.StatusBadRequest, "password must be at least 8 characters long")
	}
	if len(password) > 100 {
		return fiber.NewError(fiber.StatusBadRequest, "password must be less than 100 characters")
	}

	var hasUpper, hasLower, hasNumber, hasSpecial bool

	for _, char := range password {
		switch {
		case 'A' <= char && char <= 'Z':
			hasUpper = true
		case 'a' <= char && char <= 'z':
			hasLower = true
		case '0' <= char && char <= '9':
			hasNumber = true
		case char == '!' || char == '@' || char == '#' || char == '$' || char == '%' || char == '^' || char == '&' || char == '*':
			hasSpecial = true
		}
	}

	if !hasUpper {
		return fiber.NewError(fiber.StatusBadRequest, "password must contain at least one uppercase letter")
	}
	if !hasLower {
		return fiber.NewError(fiber.StatusBadRequest, "password must contain at least one lowercase letter")
	}
	if !hasNumber {
		return fiber.NewError(fiber.StatusBadRequest, "password must contain at least one number")
	}
	if !hasSpecial {
		return fiber.NewError(fiber.StatusBadRequest, "password must contain at least one special character (!@#$%^&*)")
	}

	return nil
}

func (h *UserHandler) sendValidationErrorResponse(c *fiber.Ctx, err error) error {
	validationErrors := make(map[string]string)

	if validationErr, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationErr {
			fieldName := fieldError.Field()
			switch fieldError.Tag() {
			case "required":
				validationErrors[fieldName] = fieldName + " is required"
			case "email":
				validationErrors[fieldName] = "Invalid email format"
			case "min":
				validationErrors[fieldName] = fieldName + " must be at least " + fieldError.Param() + " characters"
			case "max":
				validationErrors[fieldName] = fieldName + " must be less than " + fieldError.Param() + " characters"
			case "oneof":
				validationErrors[fieldName] = fieldName + " must be either 'admin' or 'user'"
			default:
				validationErrors[fieldName] = fieldName + " is invalid"
			}
		}
	}

	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"success": false,
		"message": "Validation error",
		"errors":  validationErrors,
	})
}

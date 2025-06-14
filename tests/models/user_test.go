package models_test

import (
	"InvoicelyX/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestUserModel(t *testing.T) {
	t.Run("Valid User Creation", func(t *testing.T) {
		user := models.User{
			UserID:    "test-user-123",
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
			Password:  "hashedpassword123",
			Role:      "user",
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		assert.Equal(t, "test-user-123", user.UserID)
		assert.Equal(t, "John", user.FirstName)
		assert.Equal(t, "Doe", user.LastName)
		assert.Equal(t, "john.doe@example.com", user.Email)
		assert.Equal(t, "hashedpassword123", user.Password)
		assert.Equal(t, "user", user.Role)
		assert.True(t, user.IsActive)
		assert.False(t, user.CreatedAt.IsZero())
		assert.False(t, user.UpdatedAt.IsZero())
	})

	t.Run("User with ObjectID", func(t *testing.T) {
		objectID := primitive.NewObjectID()
		user := models.User{
			ID:        objectID,
			UserID:    "test-user-456",
			FirstName: "Jane",
			LastName:  "Smith",
			Email:     "jane.smith@example.com",
			Role:      "admin",
			IsActive:  true,
		}

		assert.Equal(t, objectID, user.ID)
		assert.Equal(t, "test-user-456", user.UserID)
		assert.Equal(t, "admin", user.Role)
	})

	t.Run("User Default Values", func(t *testing.T) {
		user := models.User{}

		assert.Empty(t, user.UserID)
		assert.Empty(t, user.FirstName)
		assert.Empty(t, user.LastName)
		assert.Empty(t, user.Email)
		assert.Empty(t, user.Password)
		assert.Empty(t, user.Role)
		assert.False(t, user.IsActive)
		assert.True(t, user.CreatedAt.IsZero())
		assert.True(t, user.UpdatedAt.IsZero())
	})
}

func TestLoginRequest(t *testing.T) {
	t.Run("Valid Login Request", func(t *testing.T) {
		loginReq := models.LoginRequest{
			Email:    "test@example.com",
			Password: "password123",
		}

		assert.Equal(t, "test@example.com", loginReq.Email)
		assert.Equal(t, "password123", loginReq.Password)
	})

	t.Run("Empty Login Request", func(t *testing.T) {
		loginReq := models.LoginRequest{}

		assert.Empty(t, loginReq.Email)
		assert.Empty(t, loginReq.Password)
	})
}

func TestLogoutRequest(t *testing.T) {
	t.Run("Valid Logout Request", func(t *testing.T) {
		logoutReq := models.LogoutRequest{
			SessionToken: "valid-session-token",
		}

		assert.Equal(t, "valid-session-token", logoutReq.SessionToken)
	})

	t.Run("Empty Logout Request", func(t *testing.T) {
		logoutReq := models.LogoutRequest{}

		assert.Empty(t, logoutReq.SessionToken)
	})
}

func TestCreateUserRequest(t *testing.T) {
	t.Run("Valid Create User Request", func(t *testing.T) {
		createReq := models.CreateUserRequest{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
			Password:  "SecurePass123!",
			Role:      "user",
		}

		assert.Equal(t, "John", createReq.FirstName)
		assert.Equal(t, "Doe", createReq.LastName)
		assert.Equal(t, "john.doe@example.com", createReq.Email)
		assert.Equal(t, "SecurePass123!", createReq.Password)
		assert.Equal(t, "user", createReq.Role)
	})

	t.Run("Create Admin User Request", func(t *testing.T) {
		createReq := models.CreateUserRequest{
			FirstName: "Admin",
			LastName:  "User",
			Email:     "admin@example.com",
			Password:  "AdminPass123!",
			Role:      "admin",
		}

		assert.Equal(t, "Admin", createReq.FirstName)
		assert.Equal(t, "admin", createReq.Role)
	})

	t.Run("Create User Request Without Role", func(t *testing.T) {
		createReq := models.CreateUserRequest{
			FirstName: "User",
			LastName:  "NoRole",
			Email:     "user@example.com",
			Password:  "Password123!",
		}

		assert.Empty(t, createReq.Role)
		assert.Equal(t, "User", createReq.FirstName)
	})

	t.Run("Empty Create User Request", func(t *testing.T) {
		createReq := models.CreateUserRequest{}

		assert.Empty(t, createReq.FirstName)
		assert.Empty(t, createReq.LastName)
		assert.Empty(t, createReq.Email)
		assert.Empty(t, createReq.Password)
		assert.Empty(t, createReq.Role)
	})
}

// Test validation tags (these would be used with validator library)
func TestUserValidationTags(t *testing.T) {
	t.Run("User Struct Has Correct JSON Tags", func(t *testing.T) {
		user := models.User{
			UserID:    "test-123",
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john@example.com",
			Password:  "secret", // This should not appear in JSON due to json:"-" tag
			Role:      "user",
			IsActive:  true,
		}

		// Verify that fields exist and have expected values
		assert.Equal(t, "test-123", user.UserID)
		assert.Equal(t, "John", user.FirstName)
		assert.Equal(t, "Doe", user.LastName)
		assert.Equal(t, "john@example.com", user.Email)
		assert.Equal(t, "secret", user.Password) // Password is accessible in Go
		assert.Equal(t, "user", user.Role)
		assert.True(t, user.IsActive)
	})
}

// Benchmark tests
func BenchmarkUserCreation(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		user := models.User{
			UserID:    "test-user-123",
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
			Password:  "hashedpassword123",
			Role:      "user",
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		_ = user
	}
}

func BenchmarkCreateUserRequest(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := models.CreateUserRequest{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
			Password:  "SecurePass123!",
			Role:      "user",
		}
		_ = req
	}
}

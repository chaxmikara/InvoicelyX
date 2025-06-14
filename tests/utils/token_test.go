package utils_test

import (
	"InvoicelyX/utils"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

// TestGenerateJWT tests JWT token generation
func TestGenerateJWT(t *testing.T) {
	// Set test JWT secret
	originalSecret := os.Getenv("JWT_SECRET")
	testSecret := "test-secret-key"
	os.Setenv("JWT_SECRET", testSecret)
	defer os.Setenv("JWT_SECRET", originalSecret) // Restore original value

	tests := []struct {
		name      string
		userID    string
		email     string
		role      string
		firstName string
		lastName  string
		wantError bool
	}{
		{
			name:      "Valid user data",
			userID:    "123e4567-e89b-12d3-a456-426614174000",
			email:     "test@example.com",
			role:      "user",
			firstName: "John",
			lastName:  "Doe",
			wantError: false,
		},
		{
			name:      "Admin user",
			userID:    "123e4567-e89b-12d3-a456-426614174001",
			email:     "admin@example.com",
			role:      "admin",
			firstName: "Admin",
			lastName:  "User",
			wantError: false,
		},
		{
			name:      "Empty fields should still work",
			userID:    "",
			email:     "",
			role:      "",
			firstName: "",
			lastName:  "",
			wantError: false,
		},
		{
			name:      "Special characters in names",
			userID:    "123e4567-e89b-12d3-a456-426614174002",
			email:     "test+tag@example.com",
			role:      "user",
			firstName: "José",
			lastName:  "García-López",
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := utils.GenerateJWT(tt.userID, tt.email, tt.role, tt.firstName, tt.lastName)

			if tt.wantError {
				assert.Error(t, err)
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)

				// Verify the token can be parsed
				parsedToken, err := jwt.ParseWithClaims(token, &utils.Claims{}, func(token *jwt.Token) (interface{}, error) {
					return []byte(testSecret), nil
				})

				assert.NoError(t, err)
				assert.True(t, parsedToken.Valid)

				claims, ok := parsedToken.Claims.(*utils.Claims)
				assert.True(t, ok)
				assert.Equal(t, tt.userID, claims.UserID)
				assert.Equal(t, tt.email, claims.Email)
				assert.Equal(t, tt.role, claims.Role)
				assert.Equal(t, tt.firstName, claims.FirstName)
				assert.Equal(t, tt.lastName, claims.LastName)
				assert.Equal(t, "InvoicelyX", claims.Issuer)
				assert.Equal(t, tt.userID, claims.Subject)

				// Check expiration time (should be ~24 hours from now)
				expectedExpiry := time.Now().Add(24 * time.Hour)
				actualExpiry := claims.ExpiresAt.Time
				timeDiff := actualExpiry.Sub(expectedExpiry)
				assert.True(t, timeDiff < time.Minute && timeDiff > -time.Minute, "Expiry time should be within 1 minute of expected")
			}
		})
	}
}

// TestGenerateJWTWithoutSecret tests JWT generation with default secret
func TestGenerateJWTWithoutSecret(t *testing.T) {
	// Temporarily remove JWT_SECRET
	originalSecret := os.Getenv("JWT_SECRET")
	os.Unsetenv("JWT_SECRET")
	defer os.Setenv("JWT_SECRET", originalSecret)

	token, err := utils.GenerateJWT("test-user", "test@example.com", "user", "Test", "User")

	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Verify token was signed with default secret
	parsedToken, err := jwt.ParseWithClaims(token, &utils.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("default-secret-key"), nil
	})

	assert.NoError(t, err)
	assert.True(t, parsedToken.Valid)
}

// TestValidateJWT tests JWT token validation
func TestValidateJWT(t *testing.T) {
	testSecret := "test-validation-secret"
	os.Setenv("JWT_SECRET", testSecret)
	defer os.Unsetenv("JWT_SECRET")

	// Generate a test token
	testUserID := "test-user-123"
	testEmail := "test@example.com"
	testRole := "user"
	testFirstName := "Test"
	testLastName := "User"

	validToken, err := utils.GenerateJWT(testUserID, testEmail, testRole, testFirstName, testLastName)
	assert.NoError(t, err)

	tests := []struct {
		name      string
		token     string
		secret    string
		wantError bool
	}{
		{
			name:      "Valid token",
			token:     validToken,
			secret:    testSecret,
			wantError: false,
		},
		{
			name:      "Invalid secret",
			token:     validToken,
			secret:    "wrong-secret",
			wantError: true,
		},
		{
			name:      "Malformed token",
			token:     "invalid.token.here",
			secret:    testSecret,
			wantError: true,
		},
		{
			name:      "Empty token",
			token:     "",
			secret:    testSecret,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := utils.ValidateJWT(tt.token)

			if tt.wantError {
				assert.Error(t, err)
				assert.Nil(t, claims)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, claims)
				assert.Equal(t, testUserID, claims.UserID)
				assert.Equal(t, testEmail, claims.Email)
				assert.Equal(t, testRole, claims.Role)
				assert.Equal(t, testFirstName, claims.FirstName)
				assert.Equal(t, testLastName, claims.LastName)
			}
		})
	}
}

// TestExpiredToken tests behavior with expired tokens
func TestExpiredToken(t *testing.T) {
	testSecret := "test-expired-secret"
	os.Setenv("JWT_SECRET", testSecret)
	defer os.Unsetenv("JWT_SECRET")

	// Create an expired token
	claims := &utils.Claims{
		UserID:    "expired-user",
		Email:     "expired@example.com",
		Role:      "user",
		FirstName: "Expired",
		LastName:  "User",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)), // Expired 1 hour ago
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)), // Issued 2 hours ago
			NotBefore: jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)), // Valid from 2 hours ago
			Issuer:    "InvoicelyX",
			Subject:   "expired-user",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	expiredToken, err := token.SignedString([]byte(testSecret))
	assert.NoError(t, err)

	// Try to validate the expired token
	validatedClaims, err := utils.ValidateJWT(expiredToken)
	assert.Error(t, err)
	assert.Nil(t, validatedClaims)
	assert.Contains(t, err.Error(), "token is expired")
}

// TestTokenNotValidYet tests behavior with tokens that are not valid yet
func TestTokenNotValidYet(t *testing.T) {
	testSecret := "test-future-secret"
	os.Setenv("JWT_SECRET", testSecret)
	defer os.Unsetenv("JWT_SECRET")

	// Create a token that's not valid yet
	claims := &utils.Claims{
		UserID:    "future-user",
		Email:     "future@example.com",
		Role:      "user",
		FirstName: "Future",
		LastName:  "User",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),     // Valid for 24 hours
			IssuedAt:  jwt.NewNumericDate(time.Now()),                         // Issued now
			NotBefore: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),      // Valid from 1 hour in the future
			Issuer:    "InvoicelyX",
			Subject:   "future-user",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	futureToken, err := token.SignedString([]byte(testSecret))
	assert.NoError(t, err)

	// Try to validate the future token
	validatedClaims, err := utils.ValidateJWT(futureToken)
	assert.Error(t, err)
	assert.Nil(t, validatedClaims)
	assert.Contains(t, err.Error(), "token used before valid")
}

// Benchmark test for JWT generation
func BenchmarkGenerateJWT(b *testing.B) {
	os.Setenv("JWT_SECRET", "benchmark-secret")
	defer os.Unsetenv("JWT_SECRET")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = utils.GenerateJWT("user-123", "user@example.com", "user", "John", "Doe")
	}
}

// Benchmark test for JWT validation
func BenchmarkValidateJWT(b *testing.B) {
	testSecret := "benchmark-validation-secret"
	os.Setenv("JWT_SECRET", testSecret)
	defer os.Unsetenv("JWT_SECRET")

	// Generate a token to validate
	token, _ := utils.GenerateJWT("user-123", "user@example.com", "user", "John", "Doe")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = utils.ValidateJWT(token)
	}
}

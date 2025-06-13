package middleware

import (
    "InvoicelyX/utils"
    "strings"

    "github.com/gofiber/fiber/v2"
)

func JWTMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Get token from Authorization header
        authHeader := c.Get("Authorization")
        if authHeader == "" {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "success": false,
                "message": "Authorization header is required",
            })
        }

        // Check if header starts with "Bearer "
        if !strings.HasPrefix(authHeader, "Bearer ") {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "success": false,
                "message": "Invalid authorization header format. Use 'Bearer <token>'",
            })
        }

        // Extract token
        tokenString := strings.TrimPrefix(authHeader, "Bearer ")

        // Validate token
        claims, err := utils.ValidateJWT(tokenString)
        if err != nil {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "success": false,
                "message": "Invalid or expired token",
                "error":   err.Error(),
            })
        }

        // Store user info in context
        c.Locals("user_id", claims.UserID)
        c.Locals("email", claims.Email)
        c.Locals("role", claims.Role)
        c.Locals("first_name", claims.FirstName)
        c.Locals("last_name", claims.LastName)

        return c.Next()
    }
}

// Optional: Role-based middleware
func RequireRole(allowedRoles ...string) fiber.Handler {
    return func(c *fiber.Ctx) error {
        userRole := c.Locals("role").(string)

        for _, role := range allowedRoles {
            if userRole == role {
                return c.Next()
            }
        }

        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
            "success": false,
            "message": "Insufficient permissions",
        })
    }
}
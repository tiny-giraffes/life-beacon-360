// internal/middleware/auth.go

package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/tiny-giraffes/life-beacon-360/server/config"
)

func AuthMiddleware(c *fiber.Ctx) error {
	// Get authorization header
	authHeader := c.Get("Authorization")
	
	// Check if header is empty
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized - Missing authentication token",
		})
	}
	
	// Extract token - supports both "Bearer TOKEN" and plain "TOKEN" formats
	var token string
	if strings.HasPrefix(authHeader, "Bearer ") {
		// Extract token from "Bearer TOKEN" format
		token = strings.TrimPrefix(authHeader, "Bearer ")
	} else {
		// Use the raw header value as token
		token = authHeader
	}
	
	// Validate token
	if token == "" || token != config.AppConfig.ApiToken {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized - Invalid authentication token",
		})
	}
	
	return c.Next()
}
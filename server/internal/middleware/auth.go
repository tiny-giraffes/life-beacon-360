// internal/middleware/auth.go

package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tiny-giraffes/life-beacon-360/server/config"
)

func AuthMiddleware(c *fiber.Ctx) error {
	token := c.Get("Authorization")

	if token == "" || token != config.AppConfig.ApiToken {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized - Invalid or missing token",
		})
	}
	return c.Next()
}

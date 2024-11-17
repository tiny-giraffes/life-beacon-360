// internal/api/router.go

package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tiny-giraffes/life-beacon-360/server/internal/handlers"
	"github.com/tiny-giraffes/life-beacon-360/server/internal/middleware"
	"gorm.io/gorm"
)

// SetupRoutes sets up all the routes for the API
func SetupRoutes(app *fiber.App, db *gorm.DB) {
	api := app.Group("/api")

	// Location route
	api.Post("/locations", middleware.AuthMiddleware, handlers.CreateLocation(db))
}

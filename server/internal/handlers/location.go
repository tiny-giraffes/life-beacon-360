// internal/api/handlers/location.go

package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tiny-giraffes/life-beacon-360/server/internal/models"
	"github.com/tiny-giraffes/life-beacon-360/server/internal/repository"
	"gorm.io/gorm"
)

// CreateLocation godoc
// @Summary Save location
// @Description Stores location coordinates in the database
// @Tags Location
// @Security ApiKeyAuth  // This tells Swagger that this endpoint needs the token
// @Accept json
// @Produce json
// @Param location body models.LocationRequest true "Location data"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string "Unauthorized - Invalid or missing token"
// @Failure 500 {object} map[string]string
// @Router /api/locations [post]
func CreateLocation(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var locationReq models.LocationRequest

		// Parse JSON body into LocationRequest struct
		if err := c.BodyParser(&locationReq); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Failed to parse JSON",
			})
		}

		// Map to Location model
		location := models.Location{
			Latitude:  locationReq.Latitude,
			Longitude: locationReq.Longitude,
		}

		// Save the location using repository function
		if err := repository.SaveCoordinate(db, &location); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to save location",
			})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "Location saved successfully",
		})
	}
}

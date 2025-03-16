// Life Beacon 360
// Copyright (C) 2025 Tim Yashin/tiny-giraffes
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

// internal/api/handlers/location.go

package handlers

import (
	"fmt"

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

		fmt.Println("Received request to /api/locations")
		fmt.Printf("Headers: %v\n", c.GetReqHeaders())

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

// GetLatestLocations godoc
// @Summary Get latest locations
// @Description Retrieves the latest 10 location coordinates from the database
// @Tags Location
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {array} models.Location
// @Failure 401 {object} map[string]string "Unauthorized - Invalid or missing token"
// @Failure 500 {object} map[string]string
// @Router /api/locations [get]
func GetLatestLocations(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the latest 10 locations
		locations, err := repository.GetLatestLocations(db, 10)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to retrieve locations: " + err.Error(),
			})
		}

		// Return the locations as JSON
		return c.Status(fiber.StatusOK).JSON(locations)
	}
}

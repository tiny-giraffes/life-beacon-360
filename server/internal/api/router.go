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

// internal/api/router.go

package api

import (
	"github.com/labstack/echo/v4"
	"github.com/tiny-giraffes/life-beacon-360/server/internal/handlers"
	"github.com/tiny-giraffes/life-beacon-360/server/internal/middleware"
	"gorm.io/gorm"
)

// SetupRoutes sets up all the routes for the API
func SetupRoutes(e *echo.Echo, db *gorm.DB) {
	api := e.Group("/api")

	// Location routes
	api.POST("/locations", handlers.CreateLocation(db), middleware.AuthMiddleware)
	api.GET("/locations", handlers.GetLatestLocations(db), middleware.AuthMiddleware)
}

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

package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/tiny-giraffes/life-beacon-360/server/config"
	"github.com/tiny-giraffes/life-beacon-360/server/internal/api"
	"github.com/tiny-giraffes/life-beacon-360/server/internal/models"
	"github.com/tiny-giraffes/life-beacon-360/server/pkg/database"
)

func main() {
	// Load environment config
	config.LoadConfig()

	// Initialize the Fiber app
	app := fiber.New()

	// Connect to the database
	db, err := database.ConnectDB()
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	// Migrate the Location model
	if err := db.AutoMigrate(&models.Location{}); err != nil {
		log.Fatalf("failed to migrate coordinates table: %v", err)
	}

	// Set up API routes
	api.SetupRoutes(app, db)

	// Start the server
	log.Fatal(app.Listen(":8080"))
}

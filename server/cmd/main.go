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

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/tiny-giraffes/life-beacon-360/server/config"
	"github.com/tiny-giraffes/life-beacon-360/server/internal/api"
	"github.com/tiny-giraffes/life-beacon-360/server/internal/migrations"
	"github.com/tiny-giraffes/life-beacon-360/server/internal/models"
	"github.com/tiny-giraffes/life-beacon-360/server/pkg/database"
)

func main() {
	// Load environment config
	config.LoadConfig()

	// Initialize Echo instance
	e := echo.New()

	// CORS middleware
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.OPTIONS},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
	}))

	// Connect to the database
	db, err := database.ConnectDB()
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}
	defer func() {
		sqlDB, err := db.DB()
		if err != nil {
			log.Printf("Error getting underlying sql.DB: %v", err)
			return
		}
		if err := sqlDB.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}()

	// Run database migrations
	migrationRunner := migrations.NewMigrationRunner(db)
	if err := migrationRunner.RunMigrations(); err != nil {
		log.Fatal("Failed to run database migrations:", err)
	}

	// Set up API routes
	api.SetupRoutes(e, db)

	// Start the server
	log.Fatal(e.Start(":8080"))
}

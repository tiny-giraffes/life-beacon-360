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

package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUser     string
	DBPassword string
	DBName     string
	DBHost     string
	DBPort     string
	ApiToken   string
}

var AppConfig Config

func LoadConfig() {
	absPath, err := filepath.Abs("/home/tyashin/ml_proj/Exocortex/life-beacon-360/server/config/.env")
	if err != nil {
		panic(fmt.Sprintf("Error constructing absolute path: %v", err))
	}

	// Load .env file from the config directory
	err = godotenv.Load(absPath)
	if err != nil {
		panic(fmt.Sprintf("Error loading .env file: %v", err))
	}

	// Load environment variables into AppConfig
	AppConfig = Config{
		DBUser:     os.Getenv("POSTGRES_USER"),
		DBPassword: os.Getenv("POSTGRES_PASSWORD"),
		DBName:     os.Getenv("POSTGRES_DB"),
		DBHost:     os.Getenv("POSTGRES_HOST"),
		DBPort:     os.Getenv("POSTGRES_PORT"),
		ApiToken:   os.Getenv("API_TOKEN"), // Load the API token
	}

	// Ensure the API token is set, otherwise panic
	if AppConfig.ApiToken == "" {
		panic("API_TOKEN not set in .env file")
	}
}

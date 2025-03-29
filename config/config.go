package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds the application configuration
type Config struct {
	GATEWAY_PORT         string
	USER_SERVICE         string
	REWARD_SERVICE       string
	PAYMENT_SERVICE      string
	NOTIFICATION_SERVICE string
	MATCHING_SERVICE     string
	DOCUMENT_SERVICE     string
	SECRET_KEY           string
	ENABLE_LOGGING       bool
}

// LoadConfig loads configuration based on the environment
func LoadConfig() *Config {
	env := os.Getenv("ENV") // Set ENV=development or ENV=production

	// Choose the correct .env file
	var envFile string
	if env == "production" {
		envFile = ".env.production"
	} else {
		envFile = ".env.development"
	}

	// Load .env file
	if err := godotenv.Load(envFile); err != nil {
		log.Fatalf("[ERROR] Failed to load environment file: %v", err)
	}

	// Read values from environment variables
	return &Config{
		GATEWAY_PORT:         os.Getenv("GATEWAY_PORT"),
		USER_SERVICE:         os.Getenv("USER_SERVICE"),
		REWARD_SERVICE:       os.Getenv("REWARD_SERVICE"),
		PAYMENT_SERVICE:      os.Getenv("PAYMENT_SERVICE"),
		NOTIFICATION_SERVICE: os.Getenv("NOTIFICATION_SERVICE"),
		MATCHING_SERVICE:     os.Getenv("MATCHING_SERVICE"),
		DOCUMENT_SERVICE:     os.Getenv("DOCUMENT_SERVICE"),
		SECRET_KEY:           os.Getenv("SECRET_KEY"),
		ENABLE_LOGGING:       os.Getenv("ENABLE_LOGGING") == "true",
	}
}

package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds the application configuration
type Config struct {
	ENV                  string
	GATEWAY_PORT         string
	USER_SERVICE         string
	REWARD_SERVICE       string
	PAYMENT_SERVICE      string
	NOTIFICATION_SERVICE string
	MATCHING_SERVICE     string
	DOCUMENT_SERVICE     string
	SECRET_KEY           string
	POSTGRES_DB          string
	POSTGRES_USER        string
	POSTGRES_HOST        string
	POSTGRES_PORT        string
	POSTGRES_PASSWORD    string
	ENABLE_LOGGING       bool
	USER_TABLE           string
	DEBUG                bool
	TEST                 bool
}

func (c *Config) NewConfig() (any, any) {
	panic("unimplemented")
}

// NewConfig loads configuration based on the environment
func NewConfig() (*Config, error) {
	ENV := os.Getenv("ENV")
	if ENV == "" {
		ENV = "development"
	}

	// Load environment file based on the current environment
	envFile := ".env"
	if ENV != "production" {
		envFile = ".env." + ENV
	}

	if err := godotenv.Load(envFile); err != nil {
		log.Printf("Warning: Could not load %s file", envFile)
	}

	// Convert ENABLE_LOGGING to bool safely
	enableLogging, err := strconv.ParseBool(os.Getenv("ENABLE_LOGGING"))
	if err != nil {
		enableLogging = false // Default to false if conversion fails
	}

	config := &Config{
		ENV:                  ENV,
		GATEWAY_PORT:         getEnv("GATEWAY_PORT", "8080"),
		USER_SERVICE:         getEnv("USER_SERVICE", "http://localhost:8081"),
		REWARD_SERVICE:       getEnv("REWARD_SERVICE", "http://localhost:8082"),
		PAYMENT_SERVICE:      getEnv("PAYMENT_SERVICE", "http://localhost:8083"),
		NOTIFICATION_SERVICE: getEnv("NOTIFICATION_SERVICE", "http://localhost:8084"),
		MATCHING_SERVICE:     getEnv("MATCHING_SERVICE", "http://localhost:8085"),
		DOCUMENT_SERVICE:     getEnv("DOCUMENT_SERVICE", "http://localhost:8086"),
		SECRET_KEY:           getEnv("SECRET_KEY", "default-secret-key"),
		POSTGRES_DB:          getEnv("POSTGRES_DB", "mydatabase"),
		POSTGRES_USER:        getEnv("POSTGRES_USER", "user"),
		POSTGRES_HOST:        getEnv("POSTGRES_HOST", "localhost"),
		POSTGRES_PORT:        getEnv("POSTGRES_PORT", "5432"),
		POSTGRES_PASSWORD:    getEnv("POSTGRES_PASSWORD", "password"),
		ENABLE_LOGGING:       enableLogging,
		USER_TABLE:           "Users",
		DEBUG:                ENV != "production",
		TEST:                 ENV == "development_test" || ENV == "docker_test",
	}

	// Override USER_TABLE based on environment
	switch ENV {
	case "production_test":
		config.USER_TABLE = "ProductionTestUsers"
	case "development":
		config.USER_TABLE = "DevUsers"
	case "development_test":
		config.USER_TABLE = "TestUsers"
	case "docker", "docker_test":
		config.USER_TABLE = "DockerUsers"
	}

	return config, nil
}

// getEnv fetches an environment variable or returns a default value
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

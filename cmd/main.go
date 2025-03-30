package cmd

import (
	"log"

	"github.com/AntonyIS-chain/lost-found-gateway/config"
	app "github.com/AntonyIS-chain/lost-found-gateway/internal/adapters/app/handlers"
	"github.com/AntonyIS-chain/lost-found-gateway/internal/adapters/postgresDB"
	"github.com/AntonyIS-chain/lost-found-gateway/internal/core/services"
)

// GatewayServer initializes and starts the API Gateway
func GatewayServer() {
	// Load configuration
	conf, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database client
	dbClient, err := postgresDB.NewPostgresDBClient(conf)
	if err != nil {
		log.Fatalf("Failed to initialize database client: %v", err)
	}

	// Initialize services
	usersService := services.NewAuthenticationManagementService(dbClient, *conf)

	// Start HTTP server with initialized services
	app.InitGinRoutes(usersService, conf)
}

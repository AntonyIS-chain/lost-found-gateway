package cmd

import (
	"fmt"
	"log"
	"net/http"

	"github.com/AntonyIS-chain/lost-found-gateway/config"
	"github.com/AntonyIS-chain/lost-found-gateway/internal/adapters/app/handler"
	"github.com/AntonyIS-chain/lost-found-gateway/internal/core/services"
)

// GatewayServer initializes and starts the API Gateway
func GatewayServer() {

	conf, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	service := services.NewGatewayService()
	handler := handler.NewGatewayHandler(service)

	log.Printf("🚀 Gateway listening on :%v", conf.GATEWAY_PORT)
	http.ListenAndServe(fmt.Sprintf(":%v", conf.GATEWAY_PORT), handler)
}

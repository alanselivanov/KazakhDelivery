package main

import (
	"log"
	"os"

	"api-gateway/internal/config"
	"api-gateway/internal/routes"
)

func main() {
	cfg := config.LoadConfig()

	router := routes.SetupRouter(cfg)

	port := os.Getenv("PORT")
	if port == "" {
		port = cfg.Server.Port
	}

	log.Printf("API Gateway running on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

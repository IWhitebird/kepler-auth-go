package main

import (
	_ "kepler-auth-go/docs"
	"kepler-auth-go/internal/api"
	"kepler-auth-go/internal/config"
	"kepler-auth-go/internal/database"
	"log"
)

// @title Kepler Auth API
// @version 1.0
// @description Authentication service for Kepler platform
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.skylarklabs.ai/support
// @contact.email support@skylarklabs.ai

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8000
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	cfg := config.Load()

	if err := database.Connect(cfg); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := database.Migrate(); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	if err := database.SeedDefaultData(); err != nil {
		log.Printf("Warning: Failed to seed default data: %v", err)
	}

	server := api.NewServer(cfg)
	router := server.SetupRouter()

	log.Printf("Server starting on %s:%s", cfg.Server.Host, cfg.Server.Port)
	if err := router.Run(cfg.Server.Host + ":" + cfg.Server.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

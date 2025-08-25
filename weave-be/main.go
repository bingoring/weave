package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"weave-module/config"
	"weave-module/database"
	"weave-module/middleware"
	"weave-module/queue"
	"weave-module/redis"
	"weave-be/internal/presentation/routes"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Set Gin mode
	gin.SetMode(cfg.Server.GinMode)

	// Connect to database
	if err := database.Connect(cfg); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Run database migrations
	if err := database.Migrate(); err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
	}

	// Seed database in development mode
	if cfg.App.Environment == "development" {
		if err := database.SeedData(); err != nil {
			log.Printf("Failed to seed database: %v", err)
		}
	}

	// Connect to Redis
	if err := redis.Connect(cfg); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redis.Close()

	// Connect to message queue
	if err := queue.Connect(cfg); err != nil {
		log.Fatalf("Failed to connect to queue: %v", err)
	}
	defer queue.Close()

	// Initialize Gin router
	router := gin.New()

	// Setup middleware
	router.Use(middleware.LoggingMiddleware())
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.RequestIDMiddleware())
	router.Use(gin.Recovery())

	// Setup routes
	routes.SetupRoutes(router, cfg)

	// Start server
	port := cfg.Server.Port
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
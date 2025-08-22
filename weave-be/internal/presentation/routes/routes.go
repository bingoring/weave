package routes

import (
	"github.com/gin-gonic/gin"
	"weave-module/config"
	"weave-module/middleware"
	"weave-be/internal/application/services"
	domainServices "weave-be/internal/domain/services"
	infraDB "weave-be/internal/infrastructure/database"
	"weave-be/internal/presentation/handlers"
)

// SetupRoutes configures all API routes
// This follows the Dependency Injection pattern and Factory pattern
func SetupRoutes(router *gin.Engine, cfg *config.Config) {
	// Initialize repositories
	userRepo := infraDB.NewUserRepository()

	// Initialize domain services
	userDomainService := domainServices.NewUserDomainService(userRepo, cfg)

	// Initialize application services
	userAppService := services.NewUserApplicationService(userRepo, userDomainService)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userAppService)

	// Setup API routes
	api := router.Group("/api")
	{
		// Health check
		api.GET("/health", healthCheck)

		// Authentication routes (public)
		auth := api.Group("/auth")
		{
			auth.POST("/register", userHandler.RegisterUser)
			auth.POST("/login", userHandler.LoginUser)
		}

		// User routes
		users := api.Group("/users")
		{
			// Public routes
			users.GET("/search", userHandler.SearchUsers)
			users.GET("/:id", userHandler.GetUserByID)
			users.GET("/:id/followers", userHandler.GetFollowers)
			users.GET("/:id/following", userHandler.GetFollowing)

			// Protected routes (require authentication)
			protected := users.Group("", middleware.AuthMiddleware(cfg))
			{
				protected.GET("/profile", userHandler.GetProfile)
				protected.PUT("/profile", userHandler.UpdateProfile)
				protected.POST("/:id/follow", userHandler.FollowUser)
				protected.DELETE("/:id/follow", userHandler.UnfollowUser)
			}
		}

		// TODO: Add other resource routes
		// - Weaves routes
		// - Channels routes
		// - Comments routes
		// - Notifications routes
	}
}

// healthCheck is a simple health check endpoint
func healthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":  "ok",
		"message": "Weave API is running",
		"version": "1.0.0",
	})
}
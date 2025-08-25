package routes

import (
	"weave-be/internal/container"
	"weave-module/config"
	"weave-module/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all API routes using dependency injection container
func SetupRoutes(router *gin.Engine, cfg *config.Config) {
	// Initialize dependency injection container
	c := container.NewContainer(cfg)

	// Get handlers from container
	userHandler := c.UserHandler()
	oauthHandler := c.OAuthHandler()

	// Setup API routes
	api := router.Group("/v1/api")
	{
		// Health check
		api.GET("/health", healthCheck)

		// Authentication routes (public)
		auth := api.Group("/auth")
		{
			// Email verification authentication
			auth.POST("/send-verification", userHandler.SendEmailVerification)
			auth.POST("/verify-email", userHandler.VerifyEmailAuth)

			// OAuth routes
			auth.GET("/providers", oauthHandler.GetSupportedProviders)
			auth.GET("/google/login", oauthHandler.StartGoogleLogin)
			auth.GET("/google/callback", oauthHandler.GoogleOAuthCallback)

			// OAuth connect routes (require authentication)
			authProtected := auth.Group("", middleware.AuthMiddleware(cfg))
			{
				authProtected.GET("/google/connect", oauthHandler.StartGoogleConnect)
			}
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

		// Weave routes
		weaves := api.Group("/weaves")
		{
			// Public routes
			weaves.GET("", nil)                      // Get published weaves (with pagination)
			weaves.GET("/featured", nil)             // Get featured weaves
			weaves.GET("/trending", nil)             // Get trending weaves
			weaves.GET("/popular", nil)              // Get popular weaves
			weaves.GET("/search", nil)               // Search weaves
			weaves.GET("/:id", nil)                  // Get weave by ID
			weaves.GET("/:id/forks", nil)            // Get weave forks
			weaves.GET("/:id/versions", nil)         // Get weave versions
			weaves.GET("/:id/versions/:version", nil) // Get specific version

			// Protected routes (require authentication)
			protected := weaves.Group("", middleware.AuthMiddleware(cfg))
			{
				protected.POST("", nil)                    // Create weave
				protected.PUT("/:id", nil)                 // Update weave
				protected.DELETE("/:id", nil)              // Delete weave
				protected.POST("/:id/fork", nil)           // Fork weave
				protected.POST("/:id/like", nil)           // Like weave
				protected.DELETE("/:id/like", nil)         // Unlike weave
				protected.POST("/:id/publish", nil)        // Publish weave
				protected.POST("/:id/unpublish", nil)      // Unpublish weave
				protected.GET("/drafts", nil)              // Get user's drafts
				protected.GET("/liked", nil)               // Get liked weaves
			}
		}

		// Channel routes
		channels := api.Group("/channels")
		{
			// Public routes
			channels.GET("", nil)           // Get public channels
			channels.GET("/:id", nil)       // Get channel by ID
			channels.GET("/:id/weaves", nil) // Get weaves in channel

			// Protected routes
			protected := channels.Group("", middleware.AuthMiddleware(cfg))
			{
				protected.POST("", nil)         // Create channel
				protected.PUT("/:id", nil)      // Update channel
				protected.DELETE("/:id", nil)   // Delete channel
				protected.POST("/:id/join", nil) // Join channel
				protected.DELETE("/:id/leave", nil) // Leave channel
			}
		}

		// Collaboration routes
		collaborations := api.Group("/collaborations")
		{
			protected := collaborations.Group("", middleware.AuthMiddleware(cfg))
			{
				protected.POST("/weaves/:id/contributions", nil)     // Create contribution
				protected.GET("/weaves/:id/contributions", nil)      // Get contributions for weave
				protected.PUT("/contributions/:id", nil)             // Update contribution
				protected.DELETE("/contributions/:id", nil)          // Delete contribution
				protected.POST("/contributions/:id/review", nil)     // Review contribution
				protected.POST("/contributions/:id/merge", nil)      // Merge contribution
				protected.POST("/weaves/:id/comments", nil)          // Add comment to weave
				protected.GET("/weaves/:id/comments", nil)           // Get comments for weave
			}
		}

		// Analytics routes
		analytics := api.Group("/analytics")
		{
			protected := analytics.Group("", middleware.AuthMiddleware(cfg))
			{
				protected.GET("/dashboard", nil)        // User analytics dashboard
				protected.GET("/weaves/:id/stats", nil) // Weave statistics
			}
		}
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

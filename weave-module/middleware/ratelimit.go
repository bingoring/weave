package middleware

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"weave-module/errors"
	"weave-module/redis"
	"weave-module/utils"
)

type RateLimitConfig struct {
	MaxRequests int           // Maximum number of requests
	Window      time.Duration // Time window
	KeyFunc     func(*gin.Context) string // Function to generate rate limit key
}

func RateLimitMiddleware(config RateLimitConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := config.KeyFunc(c)
		allowed, err := checkRateLimit(c, key, config.MaxRequests, config.Window)
		if err != nil {
			utils.ErrorResponse(c, errors.InternalServerError("Rate limit check failed"))
			c.Abort()
			return
		}

		if !allowed {
			utils.ErrorResponse(c, errors.NewAppError(429, "Rate limit exceeded", ""))
			c.Abort()
			return
		}

		c.Next()
	}
}

func checkRateLimit(c *gin.Context, key string, maxRequests int, window time.Duration) (bool, error) {
	ctx := context.Background()
	redisKey := fmt.Sprintf("rate_limit:%s", key)

	// Get current count
	current, err := redis.Get(ctx, redisKey)
	if err != nil && err.Error() != "redis: nil" {
		return false, err
	}

	currentCount := 0
	if current != "" {
		currentCount, _ = strconv.Atoi(current)
	}

	if currentCount >= maxRequests {
		return false, nil
	}

	// Increment counter
	newCount := currentCount + 1
	err = redis.Set(ctx, redisKey, newCount, window)
	if err != nil {
		return false, err
	}

	return true, nil
}

// Common rate limit key generators
func IPBasedKey(c *gin.Context) string {
	return c.ClientIP()
}

func UserBasedKey(c *gin.Context) string {
	userID, exists := c.Get("user_id")
	if !exists {
		return c.ClientIP() // Fallback to IP if no user
	}
	return fmt.Sprintf("user:%s", userID)
}

func EndpointBasedKey(c *gin.Context) string {
	return fmt.Sprintf("%s:%s", c.ClientIP(), c.FullPath())
}

// Predefined rate limit configurations
var (
	// General API rate limit: 100 requests per minute per IP
	GeneralRateLimit = RateLimitConfig{
		MaxRequests: 100,
		Window:      time.Minute,
		KeyFunc:     IPBasedKey,
	}

	// Auth endpoints: 5 requests per minute per IP
	AuthRateLimit = RateLimitConfig{
		MaxRequests: 5,
		Window:      time.Minute,
		KeyFunc:     IPBasedKey,
	}

	// User-specific actions: 50 requests per minute per user
	UserRateLimit = RateLimitConfig{
		MaxRequests: 50,
		Window:      time.Minute,
		KeyFunc:     UserBasedKey,
	}

	// Heavy operations: 10 requests per minute per user
	HeavyOperationRateLimit = RateLimitConfig{
		MaxRequests: 10,
		Window:      time.Minute,
		KeyFunc:     UserBasedKey,
	}
)
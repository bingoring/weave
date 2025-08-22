package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"weave-module/config"
)

var Client *redis.Client

func Connect(cfg *config.Config) error {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}

	Client = rdb
	return nil
}

func Close() error {
	if Client == nil {
		return nil
	}
	return Client.Close()
}

func GetClient() *redis.Client {
	return Client
}

// Cache helper functions
func Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	if Client == nil {
		return fmt.Errorf("Redis client not initialized")
	}
	return Client.Set(ctx, key, value, expiration).Err()
}

func Get(ctx context.Context, key string) (string, error) {
	if Client == nil {
		return "", fmt.Errorf("Redis client not initialized")
	}
	return Client.Get(ctx, key).Result()
}

func Del(ctx context.Context, keys ...string) error {
	if Client == nil {
		return fmt.Errorf("Redis client not initialized")
	}
	return Client.Del(ctx, keys...).Err()
}

func Exists(ctx context.Context, keys ...string) (int64, error) {
	if Client == nil {
		return 0, fmt.Errorf("Redis client not initialized")
	}
	return Client.Exists(ctx, keys...).Result()
}

func Expire(ctx context.Context, key string, expiration time.Duration) error {
	if Client == nil {
		return fmt.Errorf("Redis client not initialized")
	}
	return Client.Expire(ctx, key, expiration).Err()
}

// Session management
func SetSession(ctx context.Context, sessionID string, userID string, expiration time.Duration) error {
	key := fmt.Sprintf("session:%s", sessionID)
	return Set(ctx, key, userID, expiration)
}

func GetSession(ctx context.Context, sessionID string) (string, error) {
	key := fmt.Sprintf("session:%s", sessionID)
	return Get(ctx, key)
}

func DeleteSession(ctx context.Context, sessionID string) error {
	key := fmt.Sprintf("session:%s", sessionID)
	return Del(ctx, key)
}

// Cache keys for common data
func CacheWeave(ctx context.Context, weaveID string, data interface{}, expiration time.Duration) error {
	key := fmt.Sprintf("weave:%s", weaveID)
	return Set(ctx, key, data, expiration)
}

func GetCachedWeave(ctx context.Context, weaveID string) (string, error) {
	key := fmt.Sprintf("weave:%s", weaveID)
	return Get(ctx, key)
}

func CacheUserProfile(ctx context.Context, userID string, data interface{}, expiration time.Duration) error {
	key := fmt.Sprintf("user:profile:%s", userID)
	return Set(ctx, key, data, expiration)
}

func GetCachedUserProfile(ctx context.Context, userID string) (string, error) {
	key := fmt.Sprintf("user:profile:%s", userID)
	return Get(ctx, key)
}
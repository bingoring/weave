package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"weave-module/database"
	"weave-module/redis"
)

type CleanupService struct{}

func NewCleanupService() *CleanupService {
	return &CleanupService{}
}

// CleanupExpiredSessions removes expired session data from Redis
func (s *CleanupService) CleanupExpiredSessions(ctx context.Context) error {
	log.Println("Cleaning up expired sessions...")

	redisClient := redis.GetClient()
	if redisClient == nil {
		return fmt.Errorf("redis client not available")
	}

	// Redis automatically handles TTL expiration, but we can clean up any orphaned session data
	// Get all session keys
	keys, err := redisClient.Keys(ctx, "session:*").Result()
	if err != nil {
		return fmt.Errorf("failed to get session keys: %w", err)
	}

	cleanedCount := 0
	for _, key := range keys {
		// Check if key exists (TTL might have expired)
		exists, err := redisClient.Exists(ctx, key).Result()
		if err != nil {
			continue
		}
		
		if exists == 0 {
			cleanedCount++
		}
	}

	// Also clean up any rate limiting keys that might be expired
	rateLimitKeys, err := redisClient.Keys(ctx, "rate_limit:*").Result()
	if err == nil {
		for _, key := range rateLimitKeys {
			ttl, err := redisClient.TTL(ctx, key).Result()
			if err != nil {
				continue
			}
			if ttl < 0 { // Key exists but has no TTL or is expired
				redisClient.Del(ctx, key)
				cleanedCount++
			}
		}
	}

	log.Printf("Cleaned up %d expired session/cache entries", cleanedCount)
	return nil
}

// CleanupOldLogs removes old log entries from the database
func (s *CleanupService) CleanupOldLogs(ctx context.Context) error {
	log.Println("Cleaning up old log entries...")

	db := database.GetDB()
	
	// Clean up analytics events older than 90 days
	cutoffDate := time.Now().AddDate(0, 0, -90)
	
	result := db.WithContext(ctx).Exec(`
		DELETE FROM analytics.weave_events 
		WHERE time < ?
	`, cutoffDate)

	if result.Error != nil {
		return fmt.Errorf("failed to clean up old analytics events: %w", result.Error)
	}

	log.Printf("Cleaned up %d old analytics events", result.RowsAffected)

	// Clean up old notification records (keep only last 30 days)
	notificationCutoff := time.Now().AddDate(0, 0, -30)
	
	result = db.WithContext(ctx).Exec(`
		DELETE FROM notifications 
		WHERE created_at < ? AND is_read = true
	`, notificationCutoff)

	if result.Error != nil {
		log.Printf("Failed to clean up old notifications: %v", result.Error)
	} else {
		log.Printf("Cleaned up %d old notifications", result.RowsAffected)
	}

	return nil
}

// CleanupTempFiles removes temporary files and cache entries
func (s *CleanupService) CleanupTempFiles(ctx context.Context) error {
	log.Println("Cleaning up temporary files and cache...")

	redisClient := redis.GetClient()
	if redisClient == nil {
		return fmt.Errorf("redis client not available")
	}

	// Clean up old cache entries
	cachePatterns := []string{
		"cache:weave:*",
		"cache:user:*",
		"analytics:daily_summary:*",
	}

	totalCleaned := 0
	for _, pattern := range cachePatterns {
		keys, err := redisClient.Keys(ctx, pattern).Result()
		if err != nil {
			continue
		}

		for _, key := range keys {
			// Check if cache entry is older than 7 days
			ttl, err := redisClient.TTL(ctx, key).Result()
			if err != nil {
				continue
			}
			
			// If TTL is very low or negative, delete it
			if ttl < time.Hour {
				redisClient.Del(ctx, key)
				totalCleaned++
			}
		}
	}

	// Clean up old trending data (keep only last 24 hours)
	trendingKeys, err := redisClient.Keys(ctx, "trending:*").Result()
	if err == nil {
		for _, key := range trendingKeys {
			ttl, err := redisClient.TTL(ctx, key).Result()
			if err != nil {
				continue
			}
			
			if ttl < 0 || ttl > 24*time.Hour {
				redisClient.Expire(ctx, key, 24*time.Hour)
			}
		}
	}

	log.Printf("Cleaned up %d temporary cache entries", totalCleaned)
	return nil
}

// DatabaseMaintenance performs database maintenance tasks
func (s *CleanupService) DatabaseMaintenance(ctx context.Context) error {
	log.Println("Performing database maintenance...")

	db := database.GetDB()
	
	// Update table statistics
	maintenanceQueries := []string{
		"ANALYZE users",
		"ANALYZE weaves", 
		"ANALYZE weave_likes",
		"ANALYZE user_follows",
		"ANALYZE channels",
		"ANALYZE analytics.weave_events",
	}

	for _, query := range maintenanceQueries {
		if err := db.WithContext(ctx).Exec(query).Error; err != nil {
			log.Printf("Failed to execute maintenance query '%s': %v", query, err)
			continue
		}
	}

	// Cleanup orphaned records
	cleanupQueries := []string{
		// Remove likes for deleted weaves
		`DELETE FROM weave_likes 
		 WHERE weave_id NOT IN (SELECT id FROM weaves)`,
		
		// Remove follows for deleted users
		`DELETE FROM user_follows 
		 WHERE follower_id NOT IN (SELECT id FROM users) 
		    OR following_id NOT IN (SELECT id FROM users)`,
		
		// Remove notifications for deleted users
		`DELETE FROM notifications 
		 WHERE user_id NOT IN (SELECT id FROM users)`,
	}

	totalCleaned := int64(0)
	for _, query := range cleanupQueries {
		result := db.WithContext(ctx).Exec(query)
		if result.Error != nil {
			log.Printf("Failed to execute cleanup query: %v", result.Error)
			continue
		}
		totalCleaned += result.RowsAffected
	}

	log.Printf("Database maintenance completed, cleaned %d orphaned records", totalCleaned)
	return nil
}

// ArchiveOldData moves old data to archive tables
func (s *CleanupService) ArchiveOldData(ctx context.Context) error {
	log.Println("Archiving old data...")

	db := database.GetDB()
	
	// Archive analytics events older than 1 year
	archiveDate := time.Now().AddDate(-1, 0, 0)
	
	// First, create archive table if it doesn't exist
	createArchiveTableQuery := `
		CREATE TABLE IF NOT EXISTS analytics.weave_events_archive (
			LIKE analytics.weave_events INCLUDING ALL
		)
	`
	
	if err := db.WithContext(ctx).Exec(createArchiveTableQuery).Error; err != nil {
		log.Printf("Failed to create archive table: %v", err)
		return err
	}

	// Move old data to archive
	moveToArchiveQuery := `
		INSERT INTO analytics.weave_events_archive 
		SELECT * FROM analytics.weave_events 
		WHERE time < ?
	`
	
	result := db.WithContext(ctx).Exec(moveToArchiveQuery, archiveDate)
	if result.Error != nil {
		return fmt.Errorf("failed to move data to archive: %w", result.Error)
	}

	archivedCount := result.RowsAffected

	// Delete archived data from main table
	deleteQuery := `
		DELETE FROM analytics.weave_events 
		WHERE time < ?
	`
	
	result = db.WithContext(ctx).Exec(deleteQuery, archiveDate)
	if result.Error != nil {
		log.Printf("Failed to delete archived data from main table: %v", result.Error)
	}

	// Archive old weave versions (keep only last 10 versions per weave)
	archiveVersionsQuery := `
		DELETE FROM weave_versions 
		WHERE id NOT IN (
			SELECT id FROM (
				SELECT id 
				FROM weave_versions wv1
				WHERE (
					SELECT COUNT(*) 
					FROM weave_versions wv2 
					WHERE wv2.weave_id = wv1.weave_id 
					AND wv2.version >= wv1.version
				) <= 10
			) AS keep_versions
		)
	`
	
	result = db.WithContext(ctx).Exec(archiveVersionsQuery)
	if result.Error != nil {
		log.Printf("Failed to archive old weave versions: %v", result.Error)
	} else {
		log.Printf("Archived %d old weave versions", result.RowsAffected)
	}

	log.Printf("Archived %d analytics events", archivedCount)
	return nil
}
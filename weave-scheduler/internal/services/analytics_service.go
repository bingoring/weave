package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"weave-module/database"
	"weave-module/models"
	"weave-module/redis"
)

type AnalyticsService struct{}

func NewAnalyticsService() *AnalyticsService {
	return &AnalyticsService{}
}

// ProcessAnalyticsEvents processes and aggregates analytics events
func (s *AnalyticsService) ProcessAnalyticsEvents(ctx context.Context) error {
	log.Println("Processing analytics events...")

	db := database.GetDB()
	
	// Process events from the last 5 minutes
	query := `
		INSERT INTO analytics.daily_stats (date, metric_name, metric_value, created_at)
		SELECT 
			DATE(time) as date,
			event_type as metric_name,
			COUNT(*) as metric_value,
			NOW() as created_at
		FROM analytics.weave_events 
		WHERE time >= NOW() - INTERVAL '5 minutes'
		GROUP BY DATE(time), event_type
		ON CONFLICT (date, metric_name) 
		DO UPDATE SET 
			metric_value = daily_stats.metric_value + EXCLUDED.metric_value,
			updated_at = NOW()
	`

	if err := db.WithContext(ctx).Exec(query).Error; err != nil {
		return fmt.Errorf("failed to process analytics events: %w", err)
	}

	log.Println("Analytics events processed successfully")
	return nil
}

// UpdateUserStats updates user statistics
func (s *AnalyticsService) UpdateUserStats(ctx context.Context) error {
	log.Println("Updating user statistics...")

	db := database.GetDB()
	
	// Update user profiles with latest stats
	updateQueries := []string{
		// Update followers count
		`UPDATE user_profiles SET 
			followers_count = (
				SELECT COUNT(*) FROM user_follows 
				WHERE following_id = user_profiles.user_id
			),
			updated_at = NOW()`,
		
		// Update following count
		`UPDATE user_profiles SET 
			following_count = (
				SELECT COUNT(*) FROM user_follows 
				WHERE follower_id = user_profiles.user_id
			),
			updated_at = NOW()`,
		
		// Update weaves count
		`UPDATE user_profiles SET 
			weaves_count = (
				SELECT COUNT(*) FROM weaves 
				WHERE user_id = user_profiles.user_id AND is_published = true
			),
			updated_at = NOW()`,
		
		// Update total likes received
		`UPDATE user_profiles SET 
			total_likes_received = (
				SELECT COALESCE(SUM(w.like_count), 0) FROM weaves w 
				WHERE w.user_id = user_profiles.user_id
			),
			updated_at = NOW()`,
	}

	for _, query := range updateQueries {
		if err := db.WithContext(ctx).Exec(query).Error; err != nil {
			log.Printf("Failed to execute query: %s, error: %v", query, err)
			continue
		}
	}

	log.Println("User statistics updated successfully")
	return nil
}

// UpdateWeaveStats updates weave statistics and counters
func (s *AnalyticsService) UpdateWeaveStats(ctx context.Context) error {
	log.Println("Updating weave statistics...")

	db := database.GetDB()
	
	// Update weave like counts
	updateLikesQuery := `
		UPDATE weaves SET 
			like_count = (
				SELECT COUNT(*) FROM weave_likes 
				WHERE weave_id = weaves.id
			),
			updated_at = NOW()
		WHERE updated_at < NOW() - INTERVAL '1 hour'
	`

	if err := db.WithContext(ctx).Exec(updateLikesQuery).Error; err != nil {
		log.Printf("Failed to update weave like counts: %v", err)
	}

	// Update weave comment counts (when implemented)
	updateCommentsQuery := `
		UPDATE weaves SET 
			comment_count = (
				SELECT COUNT(*) FROM lab_comments 
				WHERE weave_id = weaves.id
			),
			updated_at = NOW()
		WHERE updated_at < NOW() - INTERVAL '1 hour'
	`

	if err := db.WithContext(ctx).Exec(updateCommentsQuery).Error; err != nil {
		log.Printf("Failed to update weave comment counts: %v", err)
	}

	log.Println("Weave statistics updated successfully")
	return nil
}

// GenerateDailyReports generates daily analytics reports
func (s *AnalyticsService) GenerateDailyReports(ctx context.Context) error {
	log.Println("Generating daily analytics reports...")

	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	// Generate daily summary
	summary, err := s.generateDailySummary(ctx, yesterday)
	if err != nil {
		return fmt.Errorf("failed to generate daily summary: %w", err)
	}

	// Store summary in Redis for quick access
	key := fmt.Sprintf("analytics:daily_summary:%s", yesterday)
	if err := redis.Set(ctx, key, summary, 7*24*time.Hour); err != nil {
		log.Printf("Failed to store daily summary in Redis: %v", err)
	}

	// Generate channel-specific reports
	if err := s.generateChannelReports(ctx, yesterday); err != nil {
		log.Printf("Failed to generate channel reports: %v", err)
	}

	log.Println("Daily analytics reports generated successfully")
	return nil
}

type DailySummary struct {
	Date            string `json:"date"`
	NewUsers        int64  `json:"new_users"`
	NewWeaves       int64  `json:"new_weaves"`
	TotalViews      int64  `json:"total_views"`
	TotalLikes      int64  `json:"total_likes"`
	ActiveUsers     int64  `json:"active_users"`
	PopularChannels []struct {
		ChannelName string `json:"channel_name"`
		WeaveCount  int64  `json:"weave_count"`
	} `json:"popular_channels"`
}

func (s *AnalyticsService) generateDailySummary(ctx context.Context, date string) (*DailySummary, error) {
	db := database.GetDB()
	summary := &DailySummary{Date: date}

	// Count new users
	db.WithContext(ctx).Model(&models.User{}).
		Where("DATE(created_at) = ?", date).
		Count(&summary.NewUsers)

	// Count new weaves
	db.WithContext(ctx).Model(&models.Weave{}).
		Where("DATE(created_at) = ? AND is_published = ?", date, true).
		Count(&summary.NewWeaves)

	// Count total views from analytics events
	db.WithContext(ctx).Raw(`
		SELECT COALESCE(SUM(metric_value), 0) 
		FROM analytics.daily_stats 
		WHERE date = ? AND metric_name = 'weave_view'
	`, date).Scan(&summary.TotalViews)

	// Count total likes
	db.WithContext(ctx).Model(&models.WeaveLike{}).
		Where("DATE(created_at) = ?", date).
		Count(&summary.TotalLikes)

	// Count active users (users who performed any action)
	db.WithContext(ctx).Raw(`
		SELECT COUNT(DISTINCT user_id) 
		FROM analytics.weave_events 
		WHERE DATE(time) = ? AND user_id IS NOT NULL
	`, date).Scan(&summary.ActiveUsers)

	// Get popular channels
	db.WithContext(ctx).Raw(`
		SELECT c.name as channel_name, COUNT(w.id) as weave_count
		FROM channels c
		JOIN weaves w ON c.id = w.channel_id
		WHERE DATE(w.created_at) = ? AND w.is_published = true
		GROUP BY c.id, c.name
		ORDER BY weave_count DESC
		LIMIT 5
	`, date).Scan(&summary.PopularChannels)

	return summary, nil
}

func (s *AnalyticsService) generateChannelReports(ctx context.Context, date string) error {
	db := database.GetDB()

	// Get channel statistics
	rows, err := db.WithContext(ctx).Raw(`
		SELECT 
			c.id,
			c.name,
			COUNT(w.id) as new_weaves,
			COALESCE(SUM(w.view_count), 0) as total_views,
			COALESCE(SUM(w.like_count), 0) as total_likes
		FROM channels c
		LEFT JOIN weaves w ON c.id = w.channel_id AND DATE(w.created_at) = ?
		WHERE c.is_active = true
		GROUP BY c.id, c.name
	`, date).Rows()

	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var channelID, channelName string
		var newWeaves, totalViews, totalLikes int

		if err := rows.Scan(&channelID, &channelName, &newWeaves, &totalViews, &totalLikes); err != nil {
			continue
		}

		// Store channel report in Redis
		report := map[string]interface{}{
			"channel_id":   channelID,
			"channel_name": channelName,
			"date":         date,
			"new_weaves":   newWeaves,
			"total_views":  totalViews,
			"total_likes":  totalLikes,
		}

		key := fmt.Sprintf("analytics:channel_report:%s:%s", channelID, date)
		if err := redis.Set(ctx, key, report, 30*24*time.Hour); err != nil {
			log.Printf("Failed to store channel report for %s: %v", channelName, err)
		}
	}

	return nil
}
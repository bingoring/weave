package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"weave-module/database"
	"weave-module/models"
	"weave-module/redis"
)

type TrendsService struct{}

func NewTrendsService() *TrendsService {
	return &TrendsService{}
}

// UpdateTrendingWeaves calculates and updates trending weaves
func (s *TrendsService) UpdateTrendingWeaves(ctx context.Context) error {
	log.Println("Updating trending weaves...")
	
	// Calculate trending score based on recent activity
	// Score = (likes * 3 + views * 1 + comments * 5) / time_decay_factor
	trendingWeaves, err := s.calculateTrendingWeaves(ctx)
	if err != nil {
		return fmt.Errorf("failed to calculate trending weaves: %w", err)
	}

	// Store trending weaves in Redis
	key := "trending:weaves"
	if err := redis.Set(ctx, key, trendingWeaves, 15*time.Minute); err != nil {
		return fmt.Errorf("failed to store trending weaves: %w", err)
	}

	// Store trending weaves by channel
	for _, weave := range trendingWeaves {
		channelKey := fmt.Sprintf("trending:channel:%s", weave.ChannelID)
		
		// Get existing trending weaves for this channel
		var channelTrending []TrendingWeave
		existingData, err := redis.Get(ctx, channelKey)
		if err == nil {
			json.Unmarshal([]byte(existingData), &channelTrending)
		}

		// Add or update this weave in channel trending
		found := false
		for i, existing := range channelTrending {
			if existing.ID == weave.ID {
				channelTrending[i] = weave
				found = true
				break
			}
		}
		if !found {
			channelTrending = append(channelTrending, weave)
		}

		// Keep only top 20 trending weaves per channel
		if len(channelTrending) > 20 {
			channelTrending = channelTrending[:20]
		}

		redis.Set(ctx, channelKey, channelTrending, 15*time.Minute)
	}

	log.Printf("Updated trending weaves: %d items", len(trendingWeaves))
	return nil
}

// UpdatePopularChannels updates popular channels based on activity
func (s *TrendsService) UpdatePopularChannels(ctx context.Context) error {
	log.Println("Updating popular channels...")
	
	// Calculate channel popularity based on recent activity
	popularChannels, err := s.calculatePopularChannels(ctx)
	if err != nil {
		return fmt.Errorf("failed to calculate popular channels: %w", err)
	}

	// Store in Redis
	key := "trending:channels"
	if err := redis.Set(ctx, key, popularChannels, 1*time.Hour); err != nil {
		return fmt.Errorf("failed to store popular channels: %w", err)
	}

	log.Printf("Updated popular channels: %d items", len(popularChannels))
	return nil
}

// GenerateRecommendations generates personalized recommendations for users
func (s *TrendsService) GenerateRecommendations(ctx context.Context) error {
	log.Println("Generating user recommendations...")

	db := database.GetDB()
	
	// Get active users (users who logged in within the last 7 days)
	var activeUsers []models.User
	err := db.WithContext(ctx).
		Where("is_active = ? AND updated_at > ?", true, time.Now().AddDate(0, 0, -7)).
		Find(&activeUsers).Error
	
	if err != nil {
		return fmt.Errorf("failed to fetch active users: %w", err)
	}

	log.Printf("Generating recommendations for %d active users", len(activeUsers))

	for _, user := range activeUsers {
		recommendations, err := s.generateUserRecommendations(ctx, user.ID.String())
		if err != nil {
			log.Printf("Failed to generate recommendations for user %s: %v", user.ID, err)
			continue
		}

		// Store recommendations in Redis
		key := fmt.Sprintf("recommendations:user:%s", user.ID)
		if err := redis.Set(ctx, key, recommendations, 24*time.Hour); err != nil {
			log.Printf("Failed to store recommendations for user %s: %v", user.ID, err)
			continue
		}
	}

	log.Println("User recommendations generated successfully")
	return nil
}

type TrendingWeave struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	ChannelID    string    `json:"channel_id"`
	ChannelName  string    `json:"channel_name"`
	UserID       string    `json:"user_id"`
	Username     string    `json:"username"`
	LikeCount    int       `json:"like_count"`
	ViewCount    int       `json:"view_count"`
	CommentCount int       `json:"comment_count"`
	TrendScore   float64   `json:"trend_score"`
	CreatedAt    time.Time `json:"created_at"`
}

func (s *TrendsService) calculateTrendingWeaves(ctx context.Context) ([]TrendingWeave, error) {
	db := database.GetDB()
	
	// Get weaves from the last 7 days with their scores
	query := `
		SELECT 
			w.id,
			w.title,
			w.channel_id,
			c.name as channel_name,
			w.user_id,
			u.username,
			w.like_count,
			w.view_count,
			w.comment_count,
			w.created_at,
			-- Calculate trending score with time decay
			(
				(w.like_count * 3.0 + w.view_count * 1.0 + w.comment_count * 5.0) /
				GREATEST(1, EXTRACT(EPOCH FROM (NOW() - w.created_at)) / 3600.0)
			) as trend_score
		FROM weaves w
		JOIN channels c ON w.channel_id = c.id
		JOIN users u ON w.user_id = u.id
		WHERE w.is_published = true 
			AND w.created_at > NOW() - INTERVAL '7 days'
			AND (w.like_count > 0 OR w.view_count > 10 OR w.comment_count > 0)
		ORDER BY trend_score DESC
		LIMIT 100
	`

	rows, err := db.WithContext(ctx).Raw(query).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trendingWeaves []TrendingWeave
	for rows.Next() {
		var tw TrendingWeave
		err := rows.Scan(
			&tw.ID, &tw.Title, &tw.ChannelID, &tw.ChannelName,
			&tw.UserID, &tw.Username, &tw.LikeCount, &tw.ViewCount,
			&tw.CommentCount, &tw.CreatedAt, &tw.TrendScore,
		)
		if err != nil {
			continue
		}
		trendingWeaves = append(trendingWeaves, tw)
	}

	return trendingWeaves, nil
}

type PopularChannel struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	Slug         string  `json:"slug"`
	Description  string  `json:"description"`
	WeaveCount   int     `json:"weave_count"`
	ActivityScore float64 `json:"activity_score"`
}

func (s *TrendsService) calculatePopularChannels(ctx context.Context) ([]PopularChannel, error) {
	db := database.GetDB()
	
	query := `
		SELECT 
			c.id,
			c.name,
			c.slug,
			COALESCE(c.description, '') as description,
			COUNT(w.id) as weave_count,
			-- Activity score based on recent weaves, likes, and views
			(
				COUNT(w.id) * 10.0 +
				COALESCE(SUM(w.like_count), 0) * 2.0 +
				COALESCE(SUM(w.view_count), 0) * 0.1
			) as activity_score
		FROM channels c
		LEFT JOIN weaves w ON c.id = w.channel_id 
			AND w.is_published = true 
			AND w.created_at > NOW() - INTERVAL '7 days'
		WHERE c.is_active = true
		GROUP BY c.id, c.name, c.slug, c.description
		HAVING activity_score > 0
		ORDER BY activity_score DESC
		LIMIT 20
	`

	rows, err := db.WithContext(ctx).Raw(query).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var popularChannels []PopularChannel
	for rows.Next() {
		var pc PopularChannel
		err := rows.Scan(
			&pc.ID, &pc.Name, &pc.Slug, &pc.Description,
			&pc.WeaveCount, &pc.ActivityScore,
		)
		if err != nil {
			continue
		}
		popularChannels = append(popularChannels, pc)
	}

	return popularChannels, nil
}

type UserRecommendations struct {
	UserID               string          `json:"user_id"`
	RecommendedWeaves    []TrendingWeave `json:"recommended_weaves"`
	RecommendedChannels  []PopularChannel `json:"recommended_channels"`
	RecommendedUsers     []RecommendedUser `json:"recommended_users"`
	GeneratedAt          time.Time       `json:"generated_at"`
}

type RecommendedUser struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	WeavesCount  int    `json:"weaves_count"`
	Reason       string `json:"reason"`
}

func (s *TrendsService) generateUserRecommendations(ctx context.Context, userID string) (*UserRecommendations, error) {
	db := database.GetDB()
	
	recommendations := &UserRecommendations{
		UserID:      userID,
		GeneratedAt: time.Now(),
	}

	// Get user's channel preferences based on their activity
	var userChannels []string
	db.WithContext(ctx).Raw(`
		SELECT DISTINCT w.channel_id 
		FROM weaves w 
		WHERE w.user_id = ? 
		ORDER BY w.created_at DESC 
		LIMIT 5
	`, userID).Scan(&userChannels)

	// Recommend weaves from similar channels
	if len(userChannels) > 0 {
		channelIDs := "'" + userChannels[0] + "'"
		for i := 1; i < len(userChannels); i++ {
			channelIDs += ",'" + userChannels[i] + "'"
		}

		query := fmt.Sprintf(`
			SELECT w.id, w.title, w.channel_id, c.name as channel_name,
				   w.user_id, u.username, w.like_count, w.view_count,
				   w.comment_count, w.created_at, 
				   (w.like_count * 2.0 + w.view_count * 0.5) as score
			FROM weaves w
			JOIN channels c ON w.channel_id = c.id
			JOIN users u ON w.user_id = u.id
			WHERE w.channel_id IN (%s)
				AND w.user_id != ?
				AND w.is_published = true
				AND w.created_at > NOW() - INTERVAL '14 days'
			ORDER BY score DESC
			LIMIT 10
		`, channelIDs)

		rows, err := db.WithContext(ctx).Raw(query, userID).Rows()
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var tw TrendingWeave
				rows.Scan(&tw.ID, &tw.Title, &tw.ChannelID, &tw.ChannelName,
					&tw.UserID, &tw.Username, &tw.LikeCount, &tw.ViewCount,
					&tw.CommentCount, &tw.CreatedAt, &tw.TrendScore)
				recommendations.RecommendedWeaves = append(recommendations.RecommendedWeaves, tw)
			}
		}
	}

	// Recommend users to follow based on similar interests
	query := `
		SELECT u.id, u.username, COUNT(w.id) as weaves_count
		FROM users u
		JOIN weaves w ON u.id = w.user_id
		WHERE u.id != ?
			AND u.is_active = true
			AND w.is_published = true
			AND NOT EXISTS (
				SELECT 1 FROM user_follows uf 
				WHERE uf.follower_id = ? AND uf.following_id = u.id
			)
		GROUP BY u.id, u.username
		HAVING weaves_count >= 3
		ORDER BY weaves_count DESC
		LIMIT 5
	`

	rows, err := db.WithContext(ctx).Raw(query, userID, userID).Rows()
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var ru RecommendedUser
			rows.Scan(&ru.ID, &ru.Username, &ru.WeavesCount)
			ru.Reason = "Active creator in your areas of interest"
			recommendations.RecommendedUsers = append(recommendations.RecommendedUsers, ru)
		}
	}

	return recommendations, nil
}
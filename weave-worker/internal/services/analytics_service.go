package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"weave-module/database"
	"weave-module/queue"
)

type AnalyticsService struct{}

func NewAnalyticsService() *AnalyticsService {
	return &AnalyticsService{}
}

// ProcessEvent processes analytics events and stores them in TimescaleDB
func (s *AnalyticsService) ProcessEvent(ctx context.Context, msg queue.AnalyticsMessage) error {
	db := database.GetDB()
	
	// Parse timestamp
	eventTime, err := time.Parse(time.RFC3339, msg.Timestamp)
	if err != nil {
		log.Printf("Invalid timestamp %s, using current time: %v", msg.Timestamp, err)
		eventTime = time.Now()
	}

	// Prepare event data for insertion into TimescaleDB
	query := `
		INSERT INTO analytics.weave_events (time, weave_id, user_id, event_type, event_data, ip_address, user_agent)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	var weaveID, userID, ipAddress, userAgent interface{}
	
	// Handle optional fields
	if msg.WeaveID != "" {
		weaveID = msg.WeaveID
	}
	if msg.UserID != "" {
		userID = msg.UserID
	}
	
	// Extract additional data from event data
	if msg.Data != nil {
		if ip, ok := msg.Data.(map[string]interface{})["ip_address"]; ok {
			ipAddress = ip
		}
		if ua, ok := msg.Data.(map[string]interface{})["user_agent"]; ok {
			userAgent = ua
		}
	}

	// Execute the insert
	err = db.WithContext(ctx).Exec(
		query,
		eventTime,
		weaveID,
		userID,
		msg.Event,
		msg.Data,
		ipAddress,
		userAgent,
	).Error

	if err != nil {
		return fmt.Errorf("failed to insert analytics event: %w", err)
	}

	// Update real-time counters for specific events
	switch msg.Event {
	case "weave_view":
		s.incrementWeaveViewCount(ctx, msg.WeaveID)
	case "weave_like":
		s.incrementWeaveLikeCount(ctx, msg.WeaveID)
	case "user_login":
		s.updateUserLastSeen(ctx, msg.UserID)
	}

	return nil
}

// incrementWeaveViewCount updates the view count for a weave
func (s *AnalyticsService) incrementWeaveViewCount(ctx context.Context, weaveID string) {
	if weaveID == "" {
		return
	}

	db := database.GetDB()
	err := db.WithContext(ctx).Exec(
		"UPDATE weaves SET view_count = view_count + 1, updated_at = NOW() WHERE id = ?",
		weaveID,
	).Error

	if err != nil {
		log.Printf("Failed to increment view count for weave %s: %v", weaveID, err)
	}
}

// incrementWeaveLikeCount updates the like count for a weave
func (s *AnalyticsService) incrementWeaveLikeCount(ctx context.Context, weaveID string) {
	if weaveID == "" {
		return
	}

	db := database.GetDB()
	err := db.WithContext(ctx).Exec(
		"UPDATE weaves SET like_count = like_count + 1, updated_at = NOW() WHERE id = ?",
		weaveID,
	).Error

	if err != nil {
		log.Printf("Failed to increment like count for weave %s: %v", weaveID, err)
	}
}

// updateUserLastSeen updates the user's last seen timestamp
func (s *AnalyticsService) updateUserLastSeen(ctx context.Context, userID string) {
	if userID == "" {
		return
	}

	db := database.GetDB()
	err := db.WithContext(ctx).Exec(
		"UPDATE users SET updated_at = NOW() WHERE id = ?",
		userID,
	).Error

	if err != nil {
		log.Printf("Failed to update last seen for user %s: %v", userID, err)
	}
}
package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"weave-module/database"
	"weave-module/models"
	"weave-module/queue"
	"weave-module/redis"
)

type NotificationService struct{}

func NewNotificationService() *NotificationService {
	return &NotificationService{}
}

// SendPendingNotifications processes pending notifications and sends them
func (s *NotificationService) SendPendingNotifications(ctx context.Context) error {
	log.Println("Processing pending notifications...")

	db := database.GetDB()
	
	// Get unread notifications from the last 5 minutes
	var notifications []models.Notification
	err := db.WithContext(ctx).
		Where("created_at > ? AND created_at <= ?", 
			time.Now().Add(-5*time.Minute), 
			time.Now()).
		Find(&notifications).Error
	
	if err != nil {
		return fmt.Errorf("failed to fetch notifications: %w", err)
	}

	if len(notifications) == 0 {
		return nil
	}

	log.Printf("Found %d notifications to process", len(notifications))

	for _, notification := range notifications {
		// Send notification via queue
		msg := queue.NotificationMessage{
			UserID:  notification.UserID.String(),
			Type:    notification.Type,
			Title:   notification.Title,
			Message: notification.Message,
		}

		if notification.Data != nil {
			var data map[string]interface{}
			if err := json.Unmarshal([]byte(*notification.Data), &data); err == nil {
				msg.Data = data
			}
		}

		if err := queue.PublishNotification(msg); err != nil {
			log.Printf("Failed to publish notification %s: %v", notification.ID, err)
			continue
		}

		log.Printf("Sent notification %s to user %s", notification.ID, notification.UserID)
	}

	return nil
}

// SendDailyDigest sends daily digest emails to active users
func (s *NotificationService) SendDailyDigest(ctx context.Context) error {
	log.Println("Sending daily digest emails...")

	db := database.GetDB()
	
	// Get active users who have enabled daily notifications
	var users []models.User
	err := db.WithContext(ctx).
		Joins("JOIN notification_settings ON users.id = notification_settings.user_id").
		Where("users.is_active = ? AND notification_settings.email_contributions = ?", true, true).
		Find(&users).Error
	
	if err != nil {
		return fmt.Errorf("failed to fetch users for daily digest: %w", err)
	}

	log.Printf("Sending daily digest to %d users", len(users))

	for _, user := range users {
		// Get user's activity for the day
		digestData, err := s.getUserDailyDigest(ctx, user.ID.String())
		if err != nil {
			log.Printf("Failed to generate digest for user %s: %v", user.ID, err)
			continue
		}

		// Skip if no activity
		if digestData.TotalActivities == 0 {
			continue
		}

		// Send email
		emailMsg := queue.EmailMessage{
			To:      user.Email,
			Subject: "Your Daily Weave Digest",
			Body:    s.generateDigestHTML(digestData),
			Data: map[string]string{
				"type":    "daily_digest",
				"user_id": user.ID.String(),
			},
		}

		if err := queue.PublishEmail(emailMsg); err != nil {
			log.Printf("Failed to send daily digest to user %s: %v", user.ID, err)
			continue
		}

		log.Printf("Sent daily digest to user %s", user.Email)
	}

	return nil
}

// SendWeeklyDigest sends weekly digest emails
func (s *NotificationService) SendWeeklyDigest(ctx context.Context) error {
	log.Println("Sending weekly digest emails...")

	db := database.GetDB()
	
	// Get active users who have enabled weekly notifications
	var users []models.User
	err := db.WithContext(ctx).
		Joins("JOIN notification_settings ON users.id = notification_settings.user_id").
		Where("users.is_active = ? AND notification_settings.email_contributions = ?", true, true).
		Find(&users).Error
	
	if err != nil {
		return fmt.Errorf("failed to fetch users for weekly digest: %w", err)
	}

	log.Printf("Sending weekly digest to %d users", len(users))

	for _, user := range users {
		// Get user's activity for the week
		digestData, err := s.getUserWeeklyDigest(ctx, user.ID.String())
		if err != nil {
			log.Printf("Failed to generate weekly digest for user %s: %v", user.ID, err)
			continue
		}

		// Send email
		emailMsg := queue.EmailMessage{
			To:      user.Email,
			Subject: "Your Weekly Weave Summary",
			Body:    s.generateWeeklyDigestHTML(digestData),
			Data: map[string]string{
				"type":    "weekly_digest",
				"user_id": user.ID.String(),
			},
		}

		if err := queue.PublishEmail(emailMsg); err != nil {
			log.Printf("Failed to send weekly digest to user %s: %v", user.ID, err)
			continue
		}

		log.Printf("Sent weekly digest to user %s", user.Email)
	}

	return nil
}

// DigestData represents user activity summary
type DigestData struct {
	UserName        string
	NewWeaves       int
	NewLikes        int
	NewFollowers    int
	NewComments     int
	TrendingWeaves  []string
	TotalActivities int
}

func (s *NotificationService) getUserDailyDigest(ctx context.Context, userID string) (*DigestData, error) {
	db := database.GetDB()
	yesterday := time.Now().AddDate(0, 0, -1)
	
	var data DigestData
	
	// Count new weaves
	db.WithContext(ctx).Model(&models.Weave{}).
		Where("user_id = ? AND created_at >= ?", userID, yesterday).
		Count((*int64)(&data.NewWeaves))
	
	// Count new likes on user's weaves
	db.WithContext(ctx).Model(&models.WeaveLike{}).
		Joins("JOIN weaves ON weave_likes.weave_id = weaves.id").
		Where("weaves.user_id = ? AND weave_likes.created_at >= ?", userID, yesterday).
		Count((*int64)(&data.NewLikes))
	
	// Count new followers
	db.WithContext(ctx).Model(&models.UserFollow{}).
		Where("following_id = ? AND created_at >= ?", userID, yesterday).
		Count((*int64)(&data.NewFollowers))
	
	data.TotalActivities = data.NewWeaves + data.NewLikes + data.NewFollowers
	
	return &data, nil
}

func (s *NotificationService) getUserWeeklyDigest(ctx context.Context, userID string) (*DigestData, error) {
	db := database.GetDB()
	lastWeek := time.Now().AddDate(0, 0, -7)
	
	var data DigestData
	
	// Count activities for the week
	db.WithContext(ctx).Model(&models.Weave{}).
		Where("user_id = ? AND created_at >= ?", userID, lastWeek).
		Count((*int64)(&data.NewWeaves))
	
	db.WithContext(ctx).Model(&models.WeaveLike{}).
		Joins("JOIN weaves ON weave_likes.weave_id = weaves.id").
		Where("weaves.user_id = ? AND weave_likes.created_at >= ?", userID, lastWeek).
		Count((*int64)(&data.NewLikes))
	
	db.WithContext(ctx).Model(&models.UserFollow{}).
		Where("following_id = ? AND created_at >= ?", userID, lastWeek).
		Count((*int64)(&data.NewFollowers))
	
	data.TotalActivities = data.NewWeaves + data.NewLikes + data.NewFollowers
	
	return &data, nil
}

func (s *NotificationService) generateDigestHTML(data *DigestData) string {
	return fmt.Sprintf(`
		<h2>Your Daily Weave Digest</h2>
		<p>Here's what happened with your content today:</p>
		<ul>
			<li>New Weaves: %d</li>
			<li>New Likes: %d</li>
			<li>New Followers: %d</li>
		</ul>
		<p>Keep creating amazing content!</p>
	`, data.NewWeaves, data.NewLikes, data.NewFollowers)
}

func (s *NotificationService) generateWeeklyDigestHTML(data *DigestData) string {
	return fmt.Sprintf(`
		<h2>Your Weekly Weave Summary</h2>
		<p>Here's your activity from the past week:</p>
		<ul>
			<li>New Weaves: %d</li>
			<li>New Likes: %d</li>
			<li>New Followers: %d</li>
		</ul>
		<p>Great work this week!</p>
	`, data.NewWeaves, data.NewLikes, data.NewFollowers)
}
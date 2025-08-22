package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"weave-module/database"
	"weave-module/models"
	"weave-module/redis"
)

type NotificationService struct{}

func NewNotificationService() *NotificationService {
	return &NotificationService{}
}

// CreateNotification creates a new notification in the database
func (s *NotificationService) CreateNotification(ctx context.Context, userID, notificationType, title, message string, data map[string]interface{}) error {
	db := database.GetDB()
	
	// Marshal data to JSON
	var dataJSON *string
	if data != nil {
		dataBytes, err := json.Marshal(data)
		if err != nil {
			return fmt.Errorf("failed to marshal notification data: %w", err)
		}
		dataStr := string(dataBytes)
		dataJSON = &dataStr
	}

	// Create notification
	notification := &models.Notification{
		UserID:  models.ParseUUID(userID),
		Type:    notificationType,
		Title:   title,
		Message: message,
		Data:    dataJSON,
		IsRead:  false,
	}

	if err := db.WithContext(ctx).Create(notification).Error; err != nil {
		return fmt.Errorf("failed to create notification: %w", err)
	}

	log.Printf("Created notification for user %s: %s", userID, title)
	return nil
}

// SendPushNotification sends a push notification (placeholder for future implementation)
func (s *NotificationService) SendPushNotification(ctx context.Context, userID, title, message string, data map[string]interface{}) error {
	// TODO: Implement push notification service (Firebase, APNs, etc.)
	log.Printf("Push notification would be sent to user %s: %s - %s", userID, title, message)
	
	// For now, we'll store the notification intent in Redis for potential retry
	notificationData := map[string]interface{}{
		"user_id": userID,
		"title":   title,
		"message": message,
		"data":    data,
		"type":    "push",
	}
	
	key := fmt.Sprintf("pending_push:%s", userID)
	if err := redis.Set(ctx, key, notificationData, 0); err != nil {
		log.Printf("Failed to store pending push notification: %v", err)
	}
	
	return nil
}

// SendInAppNotification sends an in-app notification via WebSocket (placeholder)
func (s *NotificationService) SendInAppNotification(ctx context.Context, userID, title, message string, data map[string]interface{}) error {
	// TODO: Implement WebSocket notification service
	log.Printf("In-app notification would be sent to user %s: %s - %s", userID, title, message)
	
	// Store in Redis for real-time pickup by WebSocket handler
	notificationData := map[string]interface{}{
		"user_id":   userID,
		"title":     title,
		"message":   message,
		"data":      data,
		"type":      "in_app",
		"timestamp": "now",
	}
	
	key := fmt.Sprintf("realtime_notification:%s", userID)
	if err := redis.Set(ctx, key, notificationData, 300); err != nil { // 5 minute expiry
		log.Printf("Failed to store real-time notification: %v", err)
	}
	
	return nil
}

// GetUserNotificationSettings retrieves user's notification preferences
func (s *NotificationService) GetUserNotificationSettings(ctx context.Context, userID string) (*models.NotificationSetting, error) {
	db := database.GetDB()
	
	var settings models.NotificationSetting
	err := db.WithContext(ctx).Where("user_id = ?", userID).First(&settings).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get notification settings for user %s: %w", userID, err)
	}
	
	return &settings, nil
}

// ProcessNotificationByType processes different types of notifications
func (s *NotificationService) ProcessNotificationByType(ctx context.Context, notificationType, userID, title, message string, data map[string]interface{}) error {
	// Get user's notification settings
	settings, err := s.GetUserNotificationSettings(ctx, userID)
	if err != nil {
		log.Printf("Failed to get notification settings for user %s, using defaults: %v", userID, err)
		// Use default settings if not found
		settings = &models.NotificationSetting{
			EmailLikes:        true,
			EmailComments:     true,
			EmailFollows:      true,
			EmailContributions: true,
			PushLikes:         true,
			PushComments:      true,
			PushFollows:       true,
			PushContributions: true,
		}
	}

	// Always create database notification
	if err := s.CreateNotification(ctx, userID, notificationType, title, message, data); err != nil {
		log.Printf("Failed to create database notification: %v", err)
	}

	// Send push notification based on user preferences and notification type
	shouldSendPush := false
	switch notificationType {
	case "like":
		shouldSendPush = settings.PushLikes
	case "comment":
		shouldSendPush = settings.PushComments
	case "follow":
		shouldSendPush = settings.PushFollows
	case "contribution":
		shouldSendPush = settings.PushContributions
	default:
		shouldSendPush = true // Default to sending for other types
	}

	if shouldSendPush {
		if err := s.SendPushNotification(ctx, userID, title, message, data); err != nil {
			log.Printf("Failed to send push notification: %v", err)
		}
	}

	// Always send in-app notification for immediate feedback
	if err := s.SendInAppNotification(ctx, userID, title, message, data); err != nil {
		log.Printf("Failed to send in-app notification: %v", err)
	}

	return nil
}

// NotificationTemplates provides pre-defined notification messages
type NotificationTemplates struct{}

func (nt *NotificationTemplates) LikeNotification(likerUsername, weaveTitle string) (string, string) {
	title := "New Like on Your Weave"
	message := fmt.Sprintf("%s liked your weave \"%s\"", likerUsername, weaveTitle)
	return title, message
}

func (nt *NotificationTemplates) CommentNotification(commenterUsername, weaveTitle string) (string, string) {
	title := "New Comment on Your Weave"
	message := fmt.Sprintf("%s commented on your weave \"%s\"", commenterUsername, weaveTitle)
	return title, message
}

func (nt *NotificationTemplates) FollowNotification(followerUsername string) (string, string) {
	title := "New Follower"
	message := fmt.Sprintf("%s started following you", followerUsername)
	return title, message
}

func (nt *NotificationTemplates) ContributionNotification(contributorUsername, weaveTitle string) (string, string) {
	title := "New Contribution to Your Weave"
	message := fmt.Sprintf("%s made a contribution to your weave \"%s\"", contributorUsername, weaveTitle)
	return title, message
}

func (nt *NotificationTemplates) ForkNotification(forkerUsername, weaveTitle string) (string, string) {
	title := "Your Weave Was Forked"
	message := fmt.Sprintf("%s forked your weave \"%s\"", forkerUsername, weaveTitle)
	return title, message
}

func (nt *NotificationTemplates) FeatureNotification(weaveTitle string) (string, string) {
	title := "Your Weave Was Featured!"
	message := fmt.Sprintf("Congratulations! Your weave \"%s\" has been featured", weaveTitle)
	return title, message
}

// GetNotificationTemplates returns the templates instance
func GetNotificationTemplates() *NotificationTemplates {
	return &NotificationTemplates{}
}
package handlers

import (
	"context"
	"log"

	"weave-module/queue"
	"weave-worker/internal/services"
)

type NotificationHandler struct {
	notificationService *services.NotificationService
	templates           *services.NotificationTemplates
}

func NewNotificationHandler(notificationService *services.NotificationService) *NotificationHandler {
	return &NotificationHandler{
		notificationService: notificationService,
		templates:           services.GetNotificationTemplates(),
	}
}

// HandleNotification processes incoming notification messages from the queue
func (h *NotificationHandler) HandleNotification(ctx context.Context, msg queue.NotificationMessage) error {
	log.Printf("Processing notification for user %s: %s", msg.UserID, msg.Title)

	// Process the notification based on its type
	var data map[string]interface{}
	if msg.Data != nil {
		var ok bool
		data, ok = msg.Data.(map[string]interface{})
		if !ok {
			log.Printf("Invalid data format for notification: %v", msg.Data)
			data = make(map[string]interface{})
		}
	} else {
		data = make(map[string]interface{})
	}
	
	err := h.notificationService.ProcessNotificationByType(
		ctx,
		msg.Type,
		msg.UserID,
		msg.Title,
		msg.Message,
		data,
	)

	if err != nil {
		log.Printf("Failed to process notification for user %s: %v", msg.UserID, err)
		return err
	}

	log.Printf("Successfully processed notification for user %s", msg.UserID)
	return nil
}

// HandleLikeNotification handles like notifications
func (h *NotificationHandler) HandleLikeNotification(ctx context.Context, userID, likerUsername, weaveTitle, weaveID string) error {
	title, message := h.templates.LikeNotification(likerUsername, weaveTitle)
	
	data := map[string]interface{}{
		"type":       "like",
		"weave_id":   weaveID,
		"liker":      likerUsername,
		"weave_title": weaveTitle,
	}

	return h.notificationService.ProcessNotificationByType(ctx, "like", userID, title, message, data)
}

// HandleCommentNotification handles comment notifications
func (h *NotificationHandler) HandleCommentNotification(ctx context.Context, userID, commenterUsername, weaveTitle, weaveID, commentID string) error {
	title, message := h.templates.CommentNotification(commenterUsername, weaveTitle)
	
	data := map[string]interface{}{
		"type":        "comment",
		"weave_id":    weaveID,
		"comment_id":  commentID,
		"commenter":   commenterUsername,
		"weave_title": weaveTitle,
	}

	return h.notificationService.ProcessNotificationByType(ctx, "comment", userID, title, message, data)
}

// HandleFollowNotification handles follow notifications
func (h *NotificationHandler) HandleFollowNotification(ctx context.Context, userID, followerUsername, followerID string) error {
	title, message := h.templates.FollowNotification(followerUsername)
	
	data := map[string]interface{}{
		"type":        "follow",
		"follower_id": followerID,
		"follower":    followerUsername,
	}

	return h.notificationService.ProcessNotificationByType(ctx, "follow", userID, title, message, data)
}

// HandleContributionNotification handles contribution notifications
func (h *NotificationHandler) HandleContributionNotification(ctx context.Context, userID, contributorUsername, weaveTitle, weaveID, contributionID string) error {
	title, message := h.templates.ContributionNotification(contributorUsername, weaveTitle)
	
	data := map[string]interface{}{
		"type":            "contribution",
		"weave_id":        weaveID,
		"contribution_id": contributionID,
		"contributor":     contributorUsername,
		"weave_title":     weaveTitle,
	}

	return h.notificationService.ProcessNotificationByType(ctx, "contribution", userID, title, message, data)
}

// HandleForkNotification handles fork notifications
func (h *NotificationHandler) HandleForkNotification(ctx context.Context, userID, forkerUsername, weaveTitle, originalWeaveID, forkedWeaveID string) error {
	title, message := h.templates.ForkNotification(forkerUsername, weaveTitle)
	
	data := map[string]interface{}{
		"type":             "fork",
		"original_weave_id": originalWeaveID,
		"forked_weave_id":  forkedWeaveID,
		"forker":           forkerUsername,
		"weave_title":      weaveTitle,
	}

	return h.notificationService.ProcessNotificationByType(ctx, "fork", userID, title, message, data)
}

// HandleFeatureNotification handles feature notifications
func (h *NotificationHandler) HandleFeatureNotification(ctx context.Context, userID, weaveTitle, weaveID string) error {
	title, message := h.templates.FeatureNotification(weaveTitle)
	
	data := map[string]interface{}{
		"type":        "feature",
		"weave_id":    weaveID,
		"weave_title": weaveTitle,
	}

	return h.notificationService.ProcessNotificationByType(ctx, "feature", userID, title, message, data)
}

// HandleSystemNotification handles system-wide notifications
func (h *NotificationHandler) HandleSystemNotification(ctx context.Context, userID, title, message string, data map[string]interface{}) error {
	if data == nil {
		data = make(map[string]interface{})
	}
	data["type"] = "system"

	return h.notificationService.ProcessNotificationByType(ctx, "system", userID, title, message, data)
}

// HandleWelcomeNotification handles welcome notifications for new users
func (h *NotificationHandler) HandleWelcomeNotification(ctx context.Context, userID, username string) error {
	title := "Welcome to Weave!"
	message := "Start creating and collaborating with the community. Your journey begins now!"
	
	data := map[string]interface{}{
		"type":     "welcome",
		"username": username,
	}

	return h.notificationService.ProcessNotificationByType(ctx, "welcome", userID, title, message, data)
}
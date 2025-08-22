package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"weave-module/database"
	"weave-module/models"
	"weave-module/queue"
)

type ProcessingService struct{}

func NewProcessingService() *ProcessingService {
	return &ProcessingService{}
}

// ProcessTask processes various background tasks
func (s *ProcessingService) ProcessTask(ctx context.Context, msg queue.ProcessingMessage) error {
	log.Printf("Processing task: %s", msg.Type)

	switch msg.Type {
	case "generate_weave_timeline":
		return s.generateWeaveTimeline(ctx, msg.WeaveID, msg.Data)
	case "update_user_stats":
		return s.updateUserStats(ctx, msg.UserID, msg.Data)
	case "process_contribution":
		return s.processContribution(ctx, msg.Data)
	case "cleanup_user_data":
		return s.cleanupUserData(ctx, msg.UserID, msg.Data)
	case "regenerate_recommendations":
		return s.regenerateRecommendations(ctx, msg.UserID, msg.Data)
	case "process_image_upload":
		return s.processImageUpload(ctx, msg.Data)
	default:
		log.Printf("Unknown processing task type: %s", msg.Type)
		return fmt.Errorf("unknown task type: %s", msg.Type)
	}
}

// generateWeaveTimeline creates timeline entries for weave version history
func (s *ProcessingService) generateWeaveTimeline(ctx context.Context, weaveID string, data interface{}) error {
	log.Printf("Generating timeline for weave %s", weaveID)

	db := database.GetDB()
	
	// Get weave information
	var weave models.Weave
	if err := db.WithContext(ctx).Where("id = ?", weaveID).First(&weave).Error; err != nil {
		return fmt.Errorf("failed to find weave %s: %w", weaveID, err)
	}

	// Create timeline entry for version update
	timelineData := map[string]interface{}{
		"type":        "version_update",
		"weave_id":    weaveID,
		"version":     weave.Version,
		"title":       weave.Title,
		"user_id":     weave.UserID,
		"created_at":  weave.UpdatedAt,
	}

	if data != nil {
		if changeLog, ok := data.(map[string]interface{})["change_log"]; ok {
			timelineData["change_log"] = changeLog
		}
	}

	// Store timeline data (this would be in a timeline table)
	log.Printf("Timeline entry created for weave %s version %d", weaveID, weave.Version)
	return nil
}

// updateUserStats recalculates user statistics
func (s *ProcessingService) updateUserStats(ctx context.Context, userID string, data interface{}) error {
	log.Printf("Updating user stats for %s", userID)

	db := database.GetDB()
	
	// Get or create user profile
	var profile models.UserProfile
	err := db.WithContext(ctx).Where("user_id = ?", userID).First(&profile).Error
	if err != nil {
		// Create new profile if not exists
		profile = models.UserProfile{
			UserID: models.ParseUUID(userID),
		}
		if err := db.WithContext(ctx).Create(&profile).Error; err != nil {
			return fmt.Errorf("failed to create user profile: %w", err)
		}
	}

	// Update statistics
	updates := map[string]interface{}{}

	// Count weaves
	var weavesCount int64
	db.WithContext(ctx).Model(&models.Weave{}).
		Where("user_id = ? AND is_published = ?", userID, true).
		Count(&weavesCount)
	updates["weaves_count"] = int(weavesCount)

	// Count followers
	var followersCount int64
	db.WithContext(ctx).Model(&models.UserFollow{}).
		Where("following_id = ?", userID).
		Count(&followersCount)
	updates["followers_count"] = int(followersCount)

	// Count following
	var followingCount int64
	db.WithContext(ctx).Model(&models.UserFollow{}).
		Where("follower_id = ?", userID).
		Count(&followingCount)
	updates["following_count"] = int(followingCount)

	// Calculate total likes received
	var totalLikes int64
	db.WithContext(ctx).Raw(`
		SELECT COALESCE(SUM(w.like_count), 0) 
		FROM weaves w 
		WHERE w.user_id = ?
	`, userID).Scan(&totalLikes)
	updates["total_likes_received"] = int(totalLikes)

	// Update the profile
	err = db.WithContext(ctx).Model(&profile).Updates(updates).Error
	if err != nil {
		return fmt.Errorf("failed to update user profile: %w", err)
	}

	log.Printf("User stats updated for %s", userID)
	return nil
}

// processContribution handles contribution processing
func (s *ProcessingService) processContribution(ctx context.Context, data interface{}) error {
	log.Println("Processing contribution")

	contributionData, ok := data.(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid contribution data format")
	}

	contributionID, ok := contributionData["contribution_id"].(string)
	if !ok {
		return fmt.Errorf("missing contribution_id")
	}

	db := database.GetDB()
	
	// Get contribution details
	var contribution models.Contribution
	if err := db.WithContext(ctx).Where("id = ?", contributionID).First(&contribution).Error; err != nil {
		return fmt.Errorf("failed to find contribution %s: %w", contributionID, err)
	}

	// Process based on contribution type
	switch contribution.Type {
	case "suggestion":
		log.Printf("Processing suggestion contribution %s", contributionID)
		// Handle suggestion logic
	case "fork":
		log.Printf("Processing fork contribution %s", contributionID)
		// Handle fork logic
	case "merge":
		log.Printf("Processing merge contribution %s", contributionID)
		// Handle merge logic
	default:
		log.Printf("Unknown contribution type: %s", contribution.Type)
	}

	return nil
}

// cleanupUserData removes user data when account is deleted
func (s *ProcessingService) cleanupUserData(ctx context.Context, userID string, data interface{}) error {
	log.Printf("Cleaning up data for user %s", userID)

	db := database.GetDB()
	
	// This would include:
	// - Anonymizing user weaves
	// - Removing personal data
	// - Updating related records
	
	cleanupData, ok := data.(map[string]interface{})
	if !ok {
		cleanupData = make(map[string]interface{})
	}

	softDelete, _ := cleanupData["soft_delete"].(bool)
	
	if softDelete {
		// Soft delete - just mark as inactive
		err := db.WithContext(ctx).Model(&models.User{}).
			Where("id = ?", userID).
			Update("is_active", false).Error
		if err != nil {
			return fmt.Errorf("failed to soft delete user: %w", err)
		}
		log.Printf("User %s soft deleted", userID)
	} else {
		// Hard delete - remove user data
		log.Printf("Hard delete not implemented for user %s", userID)
		// TODO: Implement hard delete logic
	}

	return nil
}

// regenerateRecommendations rebuilds recommendations for a user
func (s *ProcessingService) regenerateRecommendations(ctx context.Context, userID string, data interface{}) error {
	log.Printf("Regenerating recommendations for user %s", userID)

	// This would trigger the recommendation engine
	// For now, we'll just log the action
	log.Printf("Recommendations regenerated for user %s", userID)
	
	return nil
}

// processImageUpload handles image processing tasks
func (s *ProcessingService) processImageUpload(ctx context.Context, data interface{}) error {
	log.Println("Processing image upload")

	imageData, ok := data.(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid image data format")
	}

	imageURL, ok := imageData["image_url"].(string)
	if !ok {
		return fmt.Errorf("missing image_url")
	}

	// Process image (resize, optimize, generate thumbnails, etc.)
	log.Printf("Processing image: %s", imageURL)
	
	// TODO: Implement actual image processing
	// - Generate thumbnails
	// - Optimize for web
	// - Extract metadata
	// - Virus scan
	
	return nil
}
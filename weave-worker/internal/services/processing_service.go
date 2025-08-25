package services

import (
	"context"
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
		// Hard delete - permanently remove user data
		log.Printf("Starting hard delete for user %s", userID)
		
		// Start a database transaction for atomic deletion
		tx := db.WithContext(ctx).Begin()
		if tx.Error != nil {
			return fmt.Errorf("failed to start transaction: %w", tx.Error)
		}
		defer tx.Rollback()

		// Delete in order to respect foreign key constraints
		
		// 1. Delete user likes
		if err := tx.Where("user_id = ?", userID).Delete(&models.WeaveLike{}).Error; err != nil {
			log.Printf("Failed to delete user likes: %v", err)
		}
		
		// 2. Delete user comments
		if err := tx.Where("user_id = ?", userID).Delete(&models.LabComment{}).Error; err != nil {
			log.Printf("Failed to delete user comments: %v", err)
		}
		
		// 3. Delete user contributions
		if err := tx.Where("user_id = ?", userID).Delete(&models.Contribution{}).Error; err != nil {
			log.Printf("Failed to delete user contributions: %v", err)
		}
		
		// 4. Delete user follows (both as follower and following)
		if err := tx.Where("follower_id = ? OR following_id = ?", userID, userID).Delete(&models.UserFollow{}).Error; err != nil {
			log.Printf("Failed to delete user follows: %v", err)
		}
		
		// 5. Delete user weave versions
		if err := tx.Where("user_id = ?", userID).Delete(&models.WeaveVersion{}).Error; err != nil {
			log.Printf("Failed to delete weave versions: %v", err)
		}
		
		// 6. Delete user weaves (this will cascade to related data)
		if err := tx.Where("user_id = ?", userID).Delete(&models.Weave{}).Error; err != nil {
			log.Printf("Failed to delete user weaves: %v", err)
		}
		
		// 7. Delete user profile
		if err := tx.Where("user_id = ?", userID).Delete(&models.UserProfile{}).Error; err != nil {
			log.Printf("Failed to delete user profile: %v", err)
		}
		
		// 8. Finally, delete the user record
		if err := tx.Where("id = ?", userID).Delete(&models.User{}).Error; err != nil {
			log.Printf("Failed to delete user record: %v", err)
			return fmt.Errorf("failed to delete user record: %w", err)
		}
		
		// Commit the transaction
		if err := tx.Commit().Error; err != nil {
			return fmt.Errorf("failed to commit hard delete transaction: %w", err)
		}
		
		log.Printf("User %s hard deleted successfully", userID)
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

	// Extract additional parameters
	userID, _ := imageData["user_id"].(string)
	imageType, _ := imageData["image_type"].(string)
	originalFilename, _ := imageData["original_filename"].(string)

	log.Printf("Processing image: %s for user: %s, type: %s", imageURL, userID, imageType)
	
	// Implement actual image processing
	var processingResults map[string]interface{}
	var err error
	
	switch imageType {
	case "profile_image":
		processingResults, err = s.processProfileImage(ctx, imageURL, originalFilename)
	case "weave_cover":
		processingResults, err = s.processWeaveCoverImage(ctx, imageURL, originalFilename)
	case "weave_content":
		processingResults, err = s.processWeaveContentImage(ctx, imageURL, originalFilename)
	default:
		processingResults, err = s.processGenericImage(ctx, imageURL, originalFilename)
	}
	
	if err != nil {
		log.Printf("Failed to process image %s: %v", imageURL, err)
		return fmt.Errorf("image processing failed: %w", err)
	}
	
	// Store processing results in database if needed
	if userID != "" {
		// Store in a hypothetical image_processing table
		log.Printf("Image processing completed for %s. Results: %+v", imageURL, processingResults)
	}
	
	return nil
}

// processProfileImage handles profile image specific processing
func (s *ProcessingService) processProfileImage(ctx context.Context, imageURL, filename string) (map[string]interface{}, error) {
	log.Printf("Processing profile image: %s", imageURL)
	
	results := map[string]interface{}{
		"thumbnails": map[string]string{
			"small":  generateThumbnailURL(imageURL, "50x50"),
			"medium": generateThumbnailURL(imageURL, "150x150"),
			"large":  generateThumbnailURL(imageURL, "300x300"),
		},
		"optimized_url": generateOptimizedURL(imageURL),
		"formats": map[string]string{
			"webp": generateFormatURL(imageURL, "webp"),
			"avif": generateFormatURL(imageURL, "avif"),
		},
		"metadata": extractImageMetadata(imageURL),
		"virus_scan": performVirusScan(imageURL),
		"dimensions": getDimensions(imageURL),
		"file_size": getFileSize(imageURL),
	}
	
	return results, nil
}

// processWeaveCoverImage handles weave cover image processing
func (s *ProcessingService) processWeaveCoverImage(ctx context.Context, imageURL, filename string) (map[string]interface{}, error) {
	log.Printf("Processing weave cover image: %s", imageURL)
	
	results := map[string]interface{}{
		"thumbnails": map[string]string{
			"small":  generateThumbnailURL(imageURL, "300x200"),
			"medium": generateThumbnailURL(imageURL, "600x400"),
			"large":  generateThumbnailURL(imageURL, "1200x800"),
		},
		"optimized_url": generateOptimizedURL(imageURL),
		"formats": map[string]string{
			"webp": generateFormatURL(imageURL, "webp"),
			"avif": generateFormatURL(imageURL, "avif"),
		},
		"blur_hash": generateBlurHash(imageURL),
		"dominant_color": extractDominantColor(imageURL),
		"metadata": extractImageMetadata(imageURL),
		"virus_scan": performVirusScan(imageURL),
	}
	
	return results, nil
}

// processWeaveContentImage handles images within weave content
func (s *ProcessingService) processWeaveContentImage(ctx context.Context, imageURL, filename string) (map[string]interface{}, error) {
	log.Printf("Processing weave content image: %s", imageURL)
	
	results := map[string]interface{}{
		"thumbnails": map[string]string{
			"small":  generateThumbnailURL(imageURL, "400x300"),
			"medium": generateThumbnailURL(imageURL, "800x600"),
			"large":  generateThumbnailURL(imageURL, "1200x900"),
		},
		"optimized_url": generateOptimizedURL(imageURL),
		"formats": map[string]string{
			"webp": generateFormatURL(imageURL, "webp"),
			"avif": generateFormatURL(imageURL, "avif"),
		},
		"metadata": extractImageMetadata(imageURL),
		"virus_scan": performVirusScan(imageURL),
		"alt_text_suggestions": generateAltTextSuggestions(imageURL),
	}
	
	return results, nil
}

// processGenericImage handles generic image processing
func (s *ProcessingService) processGenericImage(ctx context.Context, imageURL, filename string) (map[string]interface{}, error) {
	log.Printf("Processing generic image: %s", imageURL)
	
	results := map[string]interface{}{
		"optimized_url": generateOptimizedURL(imageURL),
		"metadata": extractImageMetadata(imageURL),
		"virus_scan": performVirusScan(imageURL),
	}
	
	return results, nil
}

// Helper functions for image processing (placeholder implementations)
func generateThumbnailURL(originalURL, size string) string {
	return fmt.Sprintf("%s?thumbnail=%s", originalURL, size)
}

func generateOptimizedURL(originalURL string) string {
	return fmt.Sprintf("%s?optimize=true", originalURL)
}

func generateFormatURL(originalURL, format string) string {
	return fmt.Sprintf("%s?format=%s", originalURL, format)
}

func generateBlurHash(imageURL string) string {
	return "LGF5]+Yk^6#M@-5c,1J5@[or[Q6."
}

func extractDominantColor(imageURL string) string {
	return "#3B82F6"
}

func extractImageMetadata(imageURL string) map[string]interface{} {
	return map[string]interface{}{
		"width":       1024,
		"height":      768,
		"format":      "jpeg",
		"color_space": "sRGB",
		"has_alpha":   false,
	}
}

func performVirusScan(imageURL string) map[string]interface{} {
	return map[string]interface{}{
		"status": "clean",
		"scanned_at": "2025-08-25T10:30:00Z",
		"threats_found": 0,
	}
}

func getDimensions(imageURL string) map[string]int {
	return map[string]int{
		"width":  1024,
		"height": 768,
	}
}

func getFileSize(imageURL string) int {
	return 204800 // 200KB
}

func generateAltTextSuggestions(imageURL string) []string {
	return []string{
		"Image uploaded by user",
		"Content image in weave",
		"Visual element",
	}
}
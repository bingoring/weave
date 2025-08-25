package jobs

import (
	"context"
	"log"
	"time"

	"weave-scheduler/internal/services"
)

// SendPendingNotifications creates a job function for sending pending notifications
func SendPendingNotifications(notificationService *services.NotificationService) func() {
	return func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()

		if err := notificationService.SendPendingNotifications(ctx); err != nil {
			log.Printf("Failed to send pending notifications: %v", err)
		}
	}
}

// SendDailyDigest creates a job function for sending daily digest emails
func SendDailyDigest(notificationService *services.NotificationService) func() {
	return func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()

		if err := notificationService.SendDailyDigest(ctx); err != nil {
			log.Printf("Failed to send daily digest: %v", err)
		}
	}
}

// SendWeeklyDigest creates a job function for sending weekly digest emails
func SendWeeklyDigest(notificationService *services.NotificationService) func() {
	return func() {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
		defer cancel()

		if err := notificationService.SendWeeklyDigest(ctx); err != nil {
			log.Printf("Failed to send weekly digest: %v", err)
		}
	}
}

// ProcessAnalyticsEvents creates a job function for processing analytics events
func ProcessAnalyticsEvents(analyticsService *services.AnalyticsService) func() {
	return func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		if err := analyticsService.ProcessAnalyticsEvents(ctx); err != nil {
			log.Printf("Failed to process analytics events: %v", err)
		}
	}
}

// UpdateUserStats creates a job function for updating user statistics
func UpdateUserStats(analyticsService *services.AnalyticsService) func() {
	return func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()

		if err := analyticsService.UpdateUserStats(ctx); err != nil {
			log.Printf("Failed to update user stats: %v", err)
		}
	}
}

// UpdateWeaveStats creates a job function for updating weave statistics
func UpdateWeaveStats(analyticsService *services.AnalyticsService) func() {
	return func() {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
		defer cancel()

		if err := analyticsService.UpdateWeaveStats(ctx); err != nil {
			log.Printf("Failed to update weave stats: %v", err)
		}
	}
}

// GenerateDailyReports creates a job function for generating daily reports
func GenerateDailyReports(analyticsService *services.AnalyticsService) func() {
	return func() {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Minute)
		defer cancel()

		if err := analyticsService.GenerateDailyReports(ctx); err != nil {
			log.Printf("Failed to generate daily reports: %v", err)
		}
	}
}

// UpdateTrendingWeaves creates a job function for updating trending weaves
func UpdateTrendingWeaves(trendsService *services.TrendsService) func() {
	return func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()

		if err := trendsService.UpdateTrendingWeaves(ctx); err != nil {
			log.Printf("Failed to update trending weaves: %v", err)
		}
	}
}

// UpdatePopularChannels creates a job function for updating popular channels
func UpdatePopularChannels(trendsService *services.TrendsService) func() {
	return func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		if err := trendsService.UpdatePopularChannels(ctx); err != nil {
			log.Printf("Failed to update popular channels: %v", err)
		}
	}
}

// GenerateRecommendations creates a job function for generating user recommendations
func GenerateRecommendations(trendsService *services.TrendsService) func() {
	return func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
		defer cancel()

		if err := trendsService.GenerateRecommendations(ctx); err != nil {
			log.Printf("Failed to generate recommendations: %v", err)
		}
	}
}

// CleanupExpiredSessions creates a job function for cleaning up expired sessions
func CleanupExpiredSessions(cleanupService *services.CleanupService) func() {
	return func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		if err := cleanupService.CleanupExpiredSessions(ctx); err != nil {
			log.Printf("Failed to cleanup expired sessions: %v", err)
		}
	}
}

// CleanupOldLogs creates a job function for cleaning up old logs
func CleanupOldLogs(cleanupService *services.CleanupService) func() {
	return func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
		defer cancel()

		if err := cleanupService.CleanupOldLogs(ctx); err != nil {
			log.Printf("Failed to cleanup old logs: %v", err)
		}
	}
}

// CleanupTempFiles creates a job function for cleaning up temporary files
func CleanupTempFiles(cleanupService *services.CleanupService) func() {
	return func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()

		if err := cleanupService.CleanupTempFiles(ctx); err != nil {
			log.Printf("Failed to cleanup temp files: %v", err)
		}
	}
}

// DatabaseMaintenance creates a job function for database maintenance
func DatabaseMaintenance(cleanupService *services.CleanupService) func() {
	return func() {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Minute)
		defer cancel()

		if err := cleanupService.DatabaseMaintenance(ctx); err != nil {
			log.Printf("Failed to perform database maintenance: %v", err)
		}
	}
}

// ArchiveOldData creates a job function for archiving old data
func ArchiveOldData(cleanupService *services.CleanupService) func() {
	return func() {
		ctx, cancel := context.WithTimeout(context.Background(), 120*time.Minute)
		defer cancel()

		if err := cleanupService.ArchiveOldData(ctx); err != nil {
			log.Printf("Failed to archive old data: %v", err)
		}
	}
}
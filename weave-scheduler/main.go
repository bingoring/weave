package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/robfig/cron/v3"
	"weave-module/config"
	"weave-module/database"
	"weave-module/queue"
	"weave-module/redis"
	"weave-scheduler/internal/jobs"
	"weave-scheduler/internal/services"
)

func main() {
	log.Println("Starting Weave Scheduler...")

	// Load configuration
	cfg := config.Load()

	// Connect to database
	if err := database.Connect(cfg); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Connect to Redis
	if err := redis.Connect(cfg); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redis.Close()

	// Connect to message queue
	if err := queue.Connect(cfg); err != nil {
		log.Fatalf("Failed to connect to queue: %v", err)
	}
	defer queue.Close()

	// Initialize services
	notificationService := services.NewNotificationService()
	analyticsService := services.NewAnalyticsService()
	cleanupService := services.NewCleanupService()
	trendsService := services.NewTrendsService()

	// Create cron scheduler with logger
	c := cron.New(cron.WithLogger(cron.VerbosePrintfLogger(log.New(os.Stdout, "cron: ", log.LstdFlags))))

	// Register jobs
	registerJobs(c, notificationService, analyticsService, cleanupService, trendsService)

	// Start the cron scheduler
	c.Start()
	log.Println("Scheduler started successfully")

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down scheduler...")

	// Stop the cron scheduler
	ctx := c.Stop()
	select {
	case <-ctx.Done():
		log.Println("Cron jobs stopped")
	case <-time.After(30 * time.Second):
		log.Println("Timeout waiting for cron jobs to stop")
	}

	log.Println("Scheduler stopped")
}

func registerJobs(
	c *cron.Cron,
	notificationService *services.NotificationService,
	analyticsService *services.AnalyticsService,
	cleanupService *services.CleanupService,
	trendsService *services.TrendsService,
) {
	// Notification jobs
	c.AddFunc("@every 1m", jobs.SendPendingNotifications(notificationService))
	c.AddFunc("0 9 * * *", jobs.SendDailyDigest(notificationService))         // Daily at 9 AM
	c.AddFunc("0 9 * * MON", jobs.SendWeeklyDigest(notificationService))      // Weekly on Monday at 9 AM

	// Analytics jobs
	c.AddFunc("@every 5m", jobs.ProcessAnalyticsEvents(analyticsService))
	c.AddFunc("@every 30m", jobs.UpdateUserStats(analyticsService))
	c.AddFunc("@every 1h", jobs.UpdateWeaveStats(analyticsService))
	c.AddFunc("0 2 * * *", jobs.GenerateDailyReports(analyticsService))       // Daily at 2 AM

	// Trending and recommendation jobs
	c.AddFunc("@every 15m", jobs.UpdateTrendingWeaves(trendsService))
	c.AddFunc("@every 1h", jobs.UpdatePopularChannels(trendsService))
	c.AddFunc("0 3 * * *", jobs.GenerateRecommendations(trendsService))       // Daily at 3 AM

	// Cleanup jobs
	c.AddFunc("0 1 * * *", jobs.CleanupExpiredSessions(cleanupService))       // Daily at 1 AM
	c.AddFunc("0 0 * * SUN", jobs.CleanupOldLogs(cleanupService))             // Weekly on Sunday at midnight
	c.AddFunc("0 4 * * *", jobs.CleanupTempFiles(cleanupService))             // Daily at 4 AM

	// Database maintenance jobs
	c.AddFunc("0 5 * * SUN", jobs.DatabaseMaintenance(cleanupService))        // Weekly on Sunday at 5 AM
	c.AddFunc("0 6 1 * *", jobs.ArchiveOldData(cleanupService))               // Monthly on 1st at 6 AM

	log.Println("All cron jobs registered successfully")
}
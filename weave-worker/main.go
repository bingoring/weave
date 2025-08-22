package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"weave-module/config"
	"weave-module/database"
	"weave-module/queue"
	"weave-module/redis"
	"weave-worker/internal/handlers"
	"weave-worker/internal/services"
)

func main() {
	log.Println("Starting Weave Worker...")

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
	emailService := services.NewEmailService(cfg)
	notificationService := services.NewNotificationService()
	analyticsService := services.NewAnalyticsService()
	processingService := services.NewProcessingService()

	// Initialize handlers
	notificationHandler := handlers.NewNotificationHandler(notificationService)
	emailHandler := handlers.NewEmailHandler(emailService)
	analyticsHandler := handlers.NewAnalyticsHandler(analyticsService)
	processingHandler := handlers.NewProcessingHandler(processingService)

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Start queue consumers
	log.Println("Starting queue consumers...")

	// Start notification consumer
	wg.Add(1)
	go func() {
		defer wg.Done()
		startNotificationConsumer(ctx, notificationHandler)
	}()

	// Start email consumer
	wg.Add(1)
	go func() {
		defer wg.Done()
		startEmailConsumer(ctx, emailHandler)
	}()

	// Start analytics consumer
	wg.Add(1)
	go func() {
		defer wg.Done()
		startAnalyticsConsumer(ctx, analyticsHandler)
	}()

	// Start processing consumer
	wg.Add(1)
	go func() {
		defer wg.Done()
		startProcessingConsumer(ctx, processingHandler)
	}()

	log.Println("Worker started successfully")

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down worker...")

	// Cancel context to stop all consumers
	cancel()

	// Wait for all goroutines to finish with timeout
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Println("All workers stopped gracefully")
	case <-time.After(30 * time.Second):
		log.Println("Timeout waiting for workers to stop")
	}

	log.Println("Worker stopped")
}

func startNotificationConsumer(ctx context.Context, handler *handlers.NotificationHandler) {
	log.Println("Starting notification consumer...")
	
	err := queue.ConsumeNotifications(func(msg queue.NotificationMessage) error {
		return handler.HandleNotification(ctx, msg)
	})
	
	if err != nil {
		log.Printf("Notification consumer error: %v", err)
	}
	
	log.Println("Notification consumer stopped")
}

func startEmailConsumer(ctx context.Context, handler *handlers.EmailHandler) {
	log.Println("Starting email consumer...")
	
	err := queue.ConsumeEmails(func(msg queue.EmailMessage) error {
		return handler.HandleEmail(ctx, msg)
	})
	
	if err != nil {
		log.Printf("Email consumer error: %v", err)
	}
	
	log.Println("Email consumer stopped")
}

func startAnalyticsConsumer(ctx context.Context, handler *handlers.AnalyticsHandler) {
	log.Println("Starting analytics consumer...")
	
	err := queue.ConsumeAnalytics(func(msg queue.AnalyticsMessage) error {
		return handler.HandleAnalytics(ctx, msg)
	})
	
	if err != nil {
		log.Printf("Analytics consumer error: %v", err)
	}
	
	log.Println("Analytics consumer stopped")
}

func startProcessingConsumer(ctx context.Context, handler *handlers.ProcessingHandler) {
	log.Println("Starting processing consumer...")
	
	err := queue.ConsumeProcessing(func(msg queue.ProcessingMessage) error {
		return handler.HandleProcessing(ctx, msg)
	})
	
	if err != nil {
		log.Printf("Processing consumer error: %v", err)
	}
	
	log.Println("Processing consumer stopped")
}
package handlers

import (
	"context"
	"log"

	"weave-module/queue"
	"weave-worker/internal/services"
)

type AnalyticsHandler struct {
	analyticsService *services.AnalyticsService
}

func NewAnalyticsHandler(analyticsService *services.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsService: analyticsService,
	}
}

// HandleAnalytics processes incoming analytics messages from the queue
func (h *AnalyticsHandler) HandleAnalytics(ctx context.Context, msg queue.AnalyticsMessage) error {
	log.Printf("Processing analytics event: %s for user %s", msg.Event, msg.UserID)

	err := h.analyticsService.ProcessEvent(ctx, msg)
	if err != nil {
		log.Printf("Failed to process analytics event %s: %v", msg.Event, err)
		return err
	}

	log.Printf("Successfully processed analytics event: %s", msg.Event)
	return nil
}
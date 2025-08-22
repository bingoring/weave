package handlers

import (
	"context"
	"log"

	"weave-module/queue"
	"weave-worker/internal/services"
)

type ProcessingHandler struct {
	processingService *services.ProcessingService
}

func NewProcessingHandler(processingService *services.ProcessingService) *ProcessingHandler {
	return &ProcessingHandler{
		processingService: processingService,
	}
}

// HandleProcessing processes incoming processing messages from the queue
func (h *ProcessingHandler) HandleProcessing(ctx context.Context, msg queue.ProcessingMessage) error {
	log.Printf("Processing task: %s for weave %s, user %s", msg.Type, msg.WeaveID, msg.UserID)

	err := h.processingService.ProcessTask(ctx, msg)
	if err != nil {
		log.Printf("Failed to process task %s: %v", msg.Type, err)
		return err
	}

	log.Printf("Successfully processed task: %s", msg.Type)
	return nil
}
package handlers

import (
	"context"
	"log"

	"weave-module/queue"
	"weave-worker/internal/services"
)

type EmailHandler struct {
	emailService *services.EmailService
}

func NewEmailHandler(emailService *services.EmailService) *EmailHandler {
	return &EmailHandler{
		emailService: emailService,
	}
}

// HandleEmail processes incoming email messages from the queue
func (h *EmailHandler) HandleEmail(ctx context.Context, msg queue.EmailMessage) error {
	log.Printf("Processing email to %s: %s", msg.To, msg.Subject)

	// Determine email type and route accordingly
	emailType := "general"
	if msg.Data != nil {
		if t, exists := msg.Data["type"]; exists {
			emailType = t
		}
	}

	var err error
	switch emailType {
	case "welcome":
		err = h.handleWelcomeEmail(msg)
	case "password_reset":
		err = h.handlePasswordResetEmail(msg)
	case "verification":
		err = h.handleVerificationEmail(msg)
	case "notification":
		err = h.handleNotificationEmail(msg)
	case "daily_digest":
		err = h.handleDailyDigestEmail(msg)
	case "weekly_digest":
		err = h.handleWeeklyDigestEmail(msg)
	default:
		err = h.handleGeneralEmail(msg)
	}

	if err != nil {
		log.Printf("Failed to send email to %s: %v", msg.To, err)
		return err
	}

	log.Printf("Successfully sent email to %s", msg.To)
	return nil
}

func (h *EmailHandler) handleWelcomeEmail(msg queue.EmailMessage) error {
	username := "User"
	if msg.Data != nil {
		if u, exists := msg.Data["username"]; exists {
			username = u
		}
	}
	
	return h.emailService.SendWelcomeEmail(msg.To, username)
}

func (h *EmailHandler) handlePasswordResetEmail(msg queue.EmailMessage) error {
	resetToken := ""
	if msg.Data != nil {
		if token, exists := msg.Data["reset_token"]; exists {
			resetToken = token
		}
	}
	
	if resetToken == "" {
		log.Printf("Missing reset token for password reset email to %s", msg.To)
		return nil // Don't retry if token is missing
	}
	
	return h.emailService.SendPasswordResetEmail(msg.To, resetToken)
}

func (h *EmailHandler) handleVerificationEmail(msg queue.EmailMessage) error {
	verificationToken := ""
	if msg.Data != nil {
		if token, exists := msg.Data["verification_token"]; exists {
			verificationToken = token
		}
	}
	
	if verificationToken == "" {
		log.Printf("Missing verification token for verification email to %s", msg.To)
		return nil // Don't retry if token is missing
	}
	
	return h.emailService.SendVerificationEmail(msg.To, verificationToken)
}

func (h *EmailHandler) handleNotificationEmail(msg queue.EmailMessage) error {
	title := msg.Subject
	if msg.Data != nil {
		if t, exists := msg.Data["title"]; exists {
			title = t
		}
	}
	
	return h.emailService.SendNotificationEmail(msg.To, title, msg.Body, msg.Data)
}

func (h *EmailHandler) handleDailyDigestEmail(msg queue.EmailMessage) error {
	// Daily digest emails come with pre-formatted HTML body
	return h.emailService.SendEmail(msg.To, msg.Subject, msg.Body, msg.Data)
}

func (h *EmailHandler) handleWeeklyDigestEmail(msg queue.EmailMessage) error {
	// Weekly digest emails come with pre-formatted HTML body
	return h.emailService.SendEmail(msg.To, msg.Subject, msg.Body, msg.Data)
}

func (h *EmailHandler) handleGeneralEmail(msg queue.EmailMessage) error {
	// Generic email handler
	return h.emailService.SendEmail(msg.To, msg.Subject, msg.Body, msg.Data)
}

// HandleBulkEmail processes bulk email campaigns (future implementation)
func (h *EmailHandler) HandleBulkEmail(ctx context.Context, recipients []string, subject, body string, data map[string]string) error {
	log.Printf("Processing bulk email to %d recipients: %s", len(recipients), subject)
	
	successCount := 0
	errorCount := 0
	
	for _, recipient := range recipients {
		err := h.emailService.SendEmail(recipient, subject, body, data)
		if err != nil {
			log.Printf("Failed to send bulk email to %s: %v", recipient, err)
			errorCount++
		} else {
			successCount++
		}
	}
	
	log.Printf("Bulk email campaign completed: %d successful, %d failed", successCount, errorCount)
	return nil
}

// HandleTransactionalEmail processes high-priority transactional emails
func (h *EmailHandler) HandleTransactionalEmail(ctx context.Context, msg queue.EmailMessage) error {
	log.Printf("Processing transactional email to %s: %s", msg.To, msg.Subject)
	
	// Transactional emails have higher priority and should be retried on failure
	maxRetries := 3
	var lastErr error
	
	for i := 0; i < maxRetries; i++ {
		err := h.emailService.SendEmail(msg.To, msg.Subject, msg.Body, msg.Data)
		if err == nil {
			log.Printf("Transactional email sent successfully to %s on attempt %d", msg.To, i+1)
			return nil
		}
		
		lastErr = err
		log.Printf("Transactional email attempt %d failed for %s: %v", i+1, msg.To, err)
	}
	
	log.Printf("Failed to send transactional email to %s after %d attempts: %v", msg.To, maxRetries, lastErr)
	return lastErr
}
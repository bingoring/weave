package entities

import (
	"time"

	"github.com/google/uuid"
)

// EmailVerification represents email verification code for authentication
type EmailVerification struct {
	ID        uuid.UUID
	Email     string
	Code      string
	ExpiresAt time.Time
	IsUsed    bool
	UserID    *uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}

// IsExpired checks if the verification code has expired
func (ev *EmailVerification) IsExpired() bool {
	return time.Now().After(ev.ExpiresAt)
}

// IsValid checks if the verification code is valid for use
func (ev *EmailVerification) IsValid() bool {
	return !ev.IsUsed && !ev.IsExpired()
}

// MarkAsUsed marks the verification code as used
func (ev *EmailVerification) MarkAsUsed(userID uuid.UUID) {
	ev.IsUsed = true
	ev.UserID = &userID
	ev.UpdatedAt = time.Now()
}

// NewEmailVerification creates a new email verification with 6-digit code
func NewEmailVerification(email, code string) *EmailVerification {
	return &EmailVerification{
		ID:        uuid.New(),
		Email:     email,
		Code:      code,
		ExpiresAt: time.Now().Add(15 * time.Minute), // 15분 후 만료
		IsUsed:    false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
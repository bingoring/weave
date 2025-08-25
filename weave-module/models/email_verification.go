package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EmailVerification struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Email     string    `gorm:"not null;index" json:"email"`
	Code      string    `gorm:"not null;unique;size:6" json:"code"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	IsUsed    bool      `gorm:"default:false" json:"is_used"`
	UserID    *uuid.UUID `gorm:"type:uuid;index" json:"user_id,omitempty"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (ev *EmailVerification) BeforeCreate(tx *gorm.DB) error {
	if ev.ID == uuid.Nil {
		ev.ID = uuid.New()
	}
	return nil
}

// Request/Response DTOs
type SendVerificationRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type VerifyEmailRequest struct {
	Code string `json:"code" binding:"required,len=6"`
}

type VerificationResponse struct {
	Message   string `json:"message"`
	ExpiresIn int    `json:"expires_in"` // seconds until expiration
}

type EmailAuthResponse struct {
	Token   string `json:"token"`
	User    User   `json:"user"`
	Message string `json:"message"`
}
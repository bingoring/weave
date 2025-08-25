package repositories

import (
	"context"

	"weave-be/internal/domain/entities"
)

// EmailVerificationRepository interface for email verification data access
type EmailVerificationRepository interface {
	// Create operations
	Create(ctx context.Context, verification *entities.EmailVerification) error
	
	// Read operations
	GetByCode(ctx context.Context, code string) (*entities.EmailVerification, error)
	GetActiveByEmail(ctx context.Context, email string) (*entities.EmailVerification, error)
	
	// Update operations
	Update(ctx context.Context, verification *entities.EmailVerification) error
	
	// Delete operations
	DeleteExpired(ctx context.Context) error
	DeleteByEmail(ctx context.Context, email string) error
	
	// Utility operations
	CountActiveByEmail(ctx context.Context, email string) (int64, error)
}
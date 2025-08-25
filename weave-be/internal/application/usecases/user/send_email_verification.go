package user

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"

	"weave-be/internal/domain/entities"
	"weave-be/internal/domain/repositories"
)

type SendEmailVerificationUseCase struct {
	emailVerificationRepo repositories.EmailVerificationRepository
	// TODO: Add email service interface when implemented
}

func NewSendEmailVerificationUseCase(
	emailVerificationRepo repositories.EmailVerificationRepository,
) *SendEmailVerificationUseCase {
	return &SendEmailVerificationUseCase{
		emailVerificationRepo: emailVerificationRepo,
	}
}

type SendEmailVerificationCommand struct {
	Email string `json:"email" validate:"required,email"`
}

type SendEmailVerificationResponse struct {
	Message   string `json:"message"`
	ExpiresIn int    `json:"expires_in"` // seconds
	Code      string `json:"code,omitempty"` // Only for development/testing
}

func (uc *SendEmailVerificationUseCase) Execute(ctx context.Context, cmd SendEmailVerificationCommand) (*SendEmailVerificationResponse, error) {
	// Generate 6-digit verification code
	code, err := uc.generateVerificationCode()
	if err != nil {
		return nil, fmt.Errorf("failed to generate verification code: %w", err)
	}

	// Delete any existing unused verification codes for this email
	if err := uc.emailVerificationRepo.DeleteByEmail(ctx, cmd.Email); err != nil {
		// Log error but don't fail the operation
		// TODO: Add proper logging
	}

	// Create new verification code
	verification := entities.NewEmailVerification(cmd.Email, code)
	if err := uc.emailVerificationRepo.Create(ctx, verification); err != nil {
		return nil, fmt.Errorf("failed to create verification code: %w", err)
	}

	// TODO: Send email with verification code
	// For now, we'll return the code for development/testing
	
	return &SendEmailVerificationResponse{
		Message:   "Verification code sent to your email",
		ExpiresIn: 900, // 15 minutes in seconds
		Code:      code, // TODO: Remove this in production
	}, nil
}

func (uc *SendEmailVerificationUseCase) generateVerificationCode() (string, error) {
	max := big.NewInt(999999)
	min := big.NewInt(100000)

	n, err := rand.Int(rand.Reader, max.Sub(max, min).Add(max, big.NewInt(1)))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%06d", n.Add(n, min).Int64()), nil
}
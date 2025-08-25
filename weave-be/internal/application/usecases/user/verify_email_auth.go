package user

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"weave-module/auth"
	"weave-module/config"
	"weave-be/internal/application/dto"
	"weave-be/internal/domain/entities"
	"weave-be/internal/domain/repositories"
)

type VerifyEmailAuthUseCase struct {
	emailVerificationRepo repositories.EmailVerificationRepository
	userRepo              repositories.UserRepository
	cfg                   *config.Config
}

func NewVerifyEmailAuthUseCase(
	emailVerificationRepo repositories.EmailVerificationRepository,
	userRepo repositories.UserRepository,
	cfg *config.Config,
) *VerifyEmailAuthUseCase {
	return &VerifyEmailAuthUseCase{
		emailVerificationRepo: emailVerificationRepo,
		userRepo:              userRepo,
		cfg:                   cfg,
	}
}

type VerifyEmailAuthCommand struct {
	Code string `json:"code" validate:"required,len=6"`
}

func (uc *VerifyEmailAuthUseCase) Execute(ctx context.Context, cmd VerifyEmailAuthCommand) (*dto.LoginResponse, error) {
	// Find verification code
	verification, err := uc.emailVerificationRepo.GetByCode(ctx, cmd.Code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid or expired verification code")
		}
		return nil, fmt.Errorf("failed to get verification code: %w", err)
	}

	// Check if code is valid (not expired and not used)
	if !verification.IsValid() {
		return nil, errors.New("verification code has expired or already been used")
	}

	// Find or create user
	var user *entities.User
	existingUser, err := uc.userRepo.GetByEmail(ctx, verification.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}

	if existingUser != nil {
		// User exists, use existing user
		user = existingUser
	} else {
		// Create new user
		username := uc.generateUsername(verification.Email)
		user = entities.NewEmailUser(username, verification.Email)
		
		if err := uc.userRepo.Create(ctx, user); err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}
	}

	// Mark verification as used
	verification.MarkAsUsed(user.ID)
	if err := uc.emailVerificationRepo.Update(ctx, verification); err != nil {
		return nil, fmt.Errorf("failed to update verification: %w", err)
	}

	// Generate JWT token
	token, err := auth.GenerateToken(user.ID, user.Username, user.Email, uc.cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &dto.LoginResponse{
		Token: token,
		User:  *dto.UserToResponse(user),
	}, nil
}

func (uc *VerifyEmailAuthUseCase) generateUsername(email string) string {
	// Extract username from email
	baseUsername := strings.Split(email, "@")[0]
	
	// Clean username (remove special characters)
	cleaned := strings.Builder{}
	for _, r := range baseUsername {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' {
			cleaned.WriteRune(r)
		}
	}
	baseUsername = strings.ToLower(cleaned.String())
	
	if len(baseUsername) < 3 {
		baseUsername = "user" + uuid.New().String()[:6]
	}
	
	// Ensure username is unique
	originalUsername := baseUsername
	counter := 1
	for {
		exists, _ := uc.userRepo.ExistsByUsername(context.Background(), baseUsername)
		if !exists {
			break
		}
		baseUsername = fmt.Sprintf("%s%d", originalUsername, counter)
		counter++
		if counter > 999 {
			baseUsername = originalUsername + uuid.New().String()[:6]
			break
		}
	}
	
	return baseUsername
}
package user

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"weave-be/internal/application/dto"
	"weave-be/internal/domain/repositories"
	"weave-module/oauth"
)

type GoogleOAuthConnectUseCase struct {
	userRepo     repositories.UserRepository
	oauthService *oauth.OAuthService
}

func NewGoogleOAuthConnectUseCase(
	userRepo repositories.UserRepository,
	oauthService *oauth.OAuthService,
) *GoogleOAuthConnectUseCase {
	return &GoogleOAuthConnectUseCase{
		userRepo:     userRepo,
		oauthService: oauthService,
	}
}

type GoogleOAuthConnectCommand struct {
	UserID uuid.UUID `json:"user_id" validate:"required"`
	Code   string    `json:"code" validate:"required"`
	State  string    `json:"state" validate:"required"`
}

func (uc *GoogleOAuthConnectUseCase) Execute(ctx context.Context, cmd GoogleOAuthConnectCommand) (*dto.UserResponse, error) {
	// Handle OAuth callback
	result, err := uc.oauthService.HandleCallback(ctx, "google", cmd.Code, cmd.State)
	if err != nil {
		return nil, err
	}

	// Verify that the state contains the correct user ID
	if result.UserID == nil || *result.UserID != cmd.UserID {
		return nil, errors.New("user ID mismatch in OAuth state")
	}

	// Get the user
	user, err := uc.userRepo.GetByID(ctx, cmd.UserID)
	if err != nil {
		return nil, err
	}

	// Check if Google account is already connected to another user
	exists, err := uc.userRepo.ExistsByGoogleID(ctx, result.Profile.ID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("Google account is already connected to another user")
	}

	// Link Google account
	user.LinkGoogleAccount(result.Profile.ID, result.Profile.Email)

	// Update user in repository
	if err := uc.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return dto.UserToResponse(user), nil
}
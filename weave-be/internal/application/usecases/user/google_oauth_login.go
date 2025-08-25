package user

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"weave-module/auth"
	"weave-module/config"
	"weave-be/internal/application/dto"
	"weave-be/internal/domain/entities"
	"weave-be/internal/domain/repositories"
	"weave-be/internal/domain/services"
	"weave-module/oauth"
)

type GoogleOAuthLoginUseCase struct {
	userRepo     repositories.UserRepository
	userDomainService services.UserDomainService
	oauthService *oauth.OAuthService
	cfg          *config.Config
}

func NewGoogleOAuthLoginUseCase(
	userRepo repositories.UserRepository,
	userDomainService services.UserDomainService,
	oauthService *oauth.OAuthService,
	cfg *config.Config,
) *GoogleOAuthLoginUseCase {
	return &GoogleOAuthLoginUseCase{
		userRepo:     userRepo,
		userDomainService: userDomainService,
		oauthService: oauthService,
		cfg:          cfg,
	}
}

type GoogleOAuthLoginCommand struct {
	Code  string `json:"code" validate:"required"`
	State string `json:"state" validate:"required"`
}

func (uc *GoogleOAuthLoginUseCase) Execute(ctx context.Context, cmd GoogleOAuthLoginCommand) (*dto.LoginResponse, error) {
	// Handle OAuth callback
	result, err := uc.oauthService.HandleCallback(ctx, "google", cmd.Code, cmd.State)
	if err != nil {
		return nil, err
	}

	// Check if user exists by Google ID
	var user *entities.User
	existingUser, err := uc.userRepo.GetByGoogleID(ctx, result.Profile.ID)
	if err != nil {
		// User doesn't exist, check if email is already taken
		emailUser, emailErr := uc.userRepo.GetByEmail(ctx, result.Profile.Email)
		if emailErr == nil {
			// Email exists but not linked to Google - link the account
			emailUser.LinkGoogleAccount(result.Profile.ID, result.Profile.Email)
			if err := uc.userRepo.Update(ctx, emailUser); err != nil {
				return nil, err
			}
			user = emailUser
		} else {
			// Create new OAuth user
			username := uc.generateUsername(result.Profile)
			user = entities.NewOAuthUser(username, result.Profile.Email, result.Profile.ID, result.Profile.Email)
			
			// Set profile image if available
			if result.Profile.Avatar != "" {
				user.ProfileImage = &result.Profile.Avatar
			}
			
			if err := uc.userRepo.Create(ctx, user); err != nil {
				return nil, err
			}
		}
	} else {
		user = existingUser
	}

	// Generate JWT token
	token, err := auth.GenerateToken(user.ID, user.Username, user.Email, uc.cfg)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Token: token,
		User:  *dto.UserToResponse(user),
	}, nil
}

func (uc *GoogleOAuthLoginUseCase) generateUsername(profile *oauth.UserProfile) string {
	baseUsername := ""
	
	if profile.DisplayName != "" {
		baseUsername = strings.ToLower(strings.ReplaceAll(profile.DisplayName, " ", ""))
	} else if profile.FirstName != "" && profile.LastName != "" {
		baseUsername = strings.ToLower(profile.FirstName + profile.LastName)
	} else if profile.FirstName != "" {
		baseUsername = strings.ToLower(profile.FirstName)
	} else {
		baseUsername = "user"
	}

	// Clean username (remove special characters)
	cleaned := strings.Builder{}
	for _, r := range baseUsername {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			cleaned.WriteRune(r)
		}
	}
	baseUsername = cleaned.String()
	
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
		baseUsername = originalUsername + string(rune('0'+counter))
		counter++
		if counter > 999 {
			baseUsername = originalUsername + uuid.New().String()[:6]
			break
		}
	}
	
	return baseUsername
}
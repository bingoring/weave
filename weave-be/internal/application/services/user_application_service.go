package services

import (
	"context"

	"github.com/google/uuid"
	"weave-module/config"
	"weave-module/errors"
	"weave-module/oauth"
	"weave-be/internal/application/commands"
	"weave-be/internal/application/dto"
	"weave-be/internal/application/queries"
	"weave-be/internal/application/usecases/user"
	"weave-be/internal/domain/repositories"
	"weave-be/internal/domain/services"
)

// UserApplicationService orchestrates user-related use cases
type UserApplicationService struct {
	// Use Cases  
	getUserProfileUC  *user.GetUserProfileUseCase
	followUserUC      *user.FollowUserUseCase
	unfollowUserUC    *user.UnfollowUserUseCase
	searchUsersUC     *user.SearchUsersUseCase
	getFollowersUC    *user.GetFollowersUseCase
	getFollowingUC    *user.GetFollowingUseCase
	
	// Email Authentication Use Cases
	sendEmailVerificationUC *user.SendEmailVerificationUseCase
	verifyEmailAuthUC       *user.VerifyEmailAuthUseCase
	
	// OAuth Use Cases
	googleOAuthLoginUC   *user.GoogleOAuthLoginUseCase
	googleOAuthConnectUC *user.GoogleOAuthConnectUseCase
	
	// Direct repository access for simple operations
	userRepo repositories.UserRepository
}

// NewUserApplicationService creates a new UserApplicationService with all use cases
func NewUserApplicationService(
	userRepo repositories.UserRepository,
	weaveRepo repositories.WeaveRepository,
	userDomainService services.UserDomainService,
	emailVerificationRepo repositories.EmailVerificationRepository,
	cfg *config.Config,
) *UserApplicationService {
	oauthService := oauth.NewOAuthService(cfg.OAuth)
	
	return &UserApplicationService{
		getUserProfileUC:  user.NewGetUserProfileUseCase(userRepo, weaveRepo),
		followUserUC:      user.NewFollowUserUseCase(userDomainService),
		unfollowUserUC:    user.NewUnfollowUserUseCase(userDomainService),
		searchUsersUC:     user.NewSearchUsersUseCase(userRepo),
		getFollowersUC:    user.NewGetFollowersUseCase(userRepo),
		getFollowingUC:    user.NewGetFollowingUseCase(userRepo),
		sendEmailVerificationUC: user.NewSendEmailVerificationUseCase(emailVerificationRepo),
		verifyEmailAuthUC:       user.NewVerifyEmailAuthUseCase(emailVerificationRepo, userRepo, cfg),
		googleOAuthLoginUC:   user.NewGoogleOAuthLoginUseCase(userRepo, userDomainService, oauthService, cfg),
		googleOAuthConnectUC: user.NewGoogleOAuthConnectUseCase(userRepo, oauthService),
		userRepo:          userRepo,
	}
}

// SendEmailVerification sends verification code to email
func (s *UserApplicationService) SendEmailVerification(ctx context.Context, email string) (*user.SendEmailVerificationResponse, error) {
	cmd := user.SendEmailVerificationCommand{
		Email: email,
	}

	return s.sendEmailVerificationUC.Execute(ctx, cmd)
}

// VerifyEmailAuth verifies email code and logs user in
func (s *UserApplicationService) VerifyEmailAuth(ctx context.Context, code string) (*dto.LoginResponse, error) {
	cmd := user.VerifyEmailAuthCommand{
		Code: code,
	}

	return s.verifyEmailAuthUC.Execute(ctx, cmd)
}

// GetUserProfile gets user profile with additional statistics
func (s *UserApplicationService) GetUserProfile(ctx context.Context, userID uuid.UUID) (*dto.UserProfileResponse, error) {
	query := queries.GetUserProfileQuery{
		UserID: userID,
	}

	return s.getUserProfileUC.Execute(ctx, query)
}

// GetUserByID gets user by ID (simple operation)
func (s *UserApplicationService) GetUserByID(ctx context.Context, userID uuid.UUID) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.NotFound("User not found")
	}

	return dto.UserToResponse(user), nil
}

// UpdateUserProfile updates user profile
func (s *UserApplicationService) UpdateUserProfile(ctx context.Context, userID uuid.UUID, req dto.UpdateUserProfileRequest) (*dto.UserResponse, error) {
	// Validate request
	if err := req.Validate(); err != nil {
		return nil, errors.BadRequestWithDetails("Invalid profile update data", err.Error())
	}

	// For now, keep this simple - could be moved to a use case later
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.NotFound("User not found")
	}

	err = s.userRepo.UpdateProfile(ctx, userID, req.ProfileImage, req.Bio)
	if err != nil {
		return nil, errors.InternalServerError("Failed to update profile")
	}

	// Get updated user
	updatedUser, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.InternalServerError("Failed to get updated user")
	}

	return dto.UserToResponse(updatedUser), nil
}

// FollowUser follows another user
func (s *UserApplicationService) FollowUser(ctx context.Context, followerID, followingID uuid.UUID) error {
	cmd := commands.FollowUserCommand{
		FollowerID:  followerID,
		FollowingID: followingID,
	}

	return s.followUserUC.Execute(ctx, cmd)
}

// UnfollowUser unfollows another user
func (s *UserApplicationService) UnfollowUser(ctx context.Context, followerID, followingID uuid.UUID) error {
	cmd := commands.UnfollowUserCommand{
		FollowerID:  followerID,
		FollowingID: followingID,
	}

	return s.unfollowUserUC.Execute(ctx, cmd)
}

// GetFollowers gets user's followers
func (s *UserApplicationService) GetFollowers(ctx context.Context, userID uuid.UUID, page, limit int) (*dto.PaginatedUsersResponse, error) {
	query := queries.GetFollowersQuery{
		UserID: userID,
		Page:   page,
		Limit:  limit,
	}

	return s.getFollowersUC.Execute(ctx, query)
}

// GetFollowing gets users that a user is following
func (s *UserApplicationService) GetFollowing(ctx context.Context, userID uuid.UUID, page, limit int) (*dto.PaginatedUsersResponse, error) {
	query := queries.GetFollowingQuery{
		UserID: userID,
		Page:   page,
		Limit:  limit,
	}

	return s.getFollowingUC.Execute(ctx, query)
}

// SearchUsers searches for users
func (s *UserApplicationService) SearchUsers(ctx context.Context, searchQuery string, page, limit int) (*dto.PaginatedUsersResponse, error) {
	query := queries.SearchUsersQuery{
		Query: searchQuery,
		Page:  page,
		Limit: limit,
	}

	return s.searchUsersUC.Execute(ctx, query)
}

// Google OAuth operations
func (s *UserApplicationService) GoogleOAuthLogin(ctx context.Context, code, state string) (*dto.LoginResponse, error) {
	cmd := user.GoogleOAuthLoginCommand{
		Code:  code,
		State: state,
	}

	return s.googleOAuthLoginUC.Execute(ctx, cmd)
}

func (s *UserApplicationService) GoogleOAuthConnect(ctx context.Context, userID uuid.UUID, code, state string) (*dto.UserResponse, error) {
	cmd := user.GoogleOAuthConnectCommand{
		UserID: userID,
		Code:   code,
		State:  state,
	}

	return s.googleOAuthConnectUC.Execute(ctx, cmd)
}
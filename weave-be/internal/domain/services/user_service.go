package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"weave-module/auth"
	"weave-module/config"
	"weave-module/errors"
	"weave-be/internal/domain/entities"
	"weave-be/internal/domain/repositories"
)

// UserDomainService handles complex business logic related to users
// This follows the Domain Service pattern for operations that don't naturally fit in a single entity
type UserDomainService interface {
	RegisterUser(ctx context.Context, username, email, password string) (*entities.User, error)
	AuthenticateUser(ctx context.Context, email, password string) (*entities.User, string, error)
	ValidateUserRegistration(ctx context.Context, username, email string) error
	CanUserCreateWeave(ctx context.Context, userID uuid.UUID) (bool, error)
	UpdateUserProfile(ctx context.Context, userID uuid.UUID, profileImage *string, bio *string) error
	FollowUser(ctx context.Context, followerID, followingID uuid.UUID) error
	UnfollowUser(ctx context.Context, followerID, followingID uuid.UUID) error
}

type userDomainService struct {
	userRepo repositories.UserRepository
	cfg      *config.Config
}

// NewUserDomainService creates a new user domain service
func NewUserDomainService(userRepo repositories.UserRepository, cfg *config.Config) UserDomainService {
	return &userDomainService{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

func (s *userDomainService) RegisterUser(ctx context.Context, username, email, password string) (*entities.User, error) {
	// Validate registration data
	if err := s.ValidateUserRegistration(ctx, username, email); err != nil {
		return nil, err
	}

	// Validate password strength
	if err := auth.ValidatePasswordStrength(password); err != nil {
		return nil, errors.BadRequest(fmt.Sprintf("Invalid password: %s", err.Error()))
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return nil, errors.InternalServerError("Failed to hash password")
	}

	// Create user entity
	user := entities.NewUser(username, email, hashedPassword)

	// Save to repository
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, errors.InternalServerError("Failed to create user")
	}

	return user, nil
}

func (s *userDomainService) AuthenticateUser(ctx context.Context, email, password string) (*entities.User, string, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, "", errors.Unauthorized("Invalid credentials")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, "", errors.Unauthorized("Account is deactivated")
	}

	// Verify password
	if err := auth.CheckPassword(password, user.PasswordHash); err != nil {
		return nil, "", errors.Unauthorized("Invalid credentials")
	}

	// Generate JWT token
	token, err := auth.GenerateToken(user.ID, user.Username, user.Email, s.cfg)
	if err != nil {
		return nil, "", errors.InternalServerError("Failed to generate token")
	}

	return user, token, nil
}

func (s *userDomainService) ValidateUserRegistration(ctx context.Context, username, email string) error {
	// Check if username already exists
	exists, err := s.userRepo.ExistsByUsername(ctx, username)
	if err != nil {
		return errors.InternalServerError("Failed to check username availability")
	}
	if exists {
		return errors.Conflict("Username already exists")
	}

	// Check if email already exists
	exists, err = s.userRepo.ExistsByEmail(ctx, email)
	if err != nil {
		return errors.InternalServerError("Failed to check email availability")
	}
	if exists {
		return errors.Conflict("Email already exists")
	}

	return nil
}

func (s *userDomainService) CanUserCreateWeave(ctx context.Context, userID uuid.UUID) (bool, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return false, err
	}

	return user.CanCreateWeave(), nil
}

func (s *userDomainService) UpdateUserProfile(ctx context.Context, userID uuid.UUID, profileImage *string, bio *string) error {
	// Verify user exists
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return errors.NotFound("User not found")
	}

	if !user.IsActive {
		return errors.Forbidden("Cannot update profile for inactive user")
	}

	return s.userRepo.UpdateProfile(ctx, userID, profileImage, bio)
}

func (s *userDomainService) FollowUser(ctx context.Context, followerID, followingID uuid.UUID) error {
	// Cannot follow yourself
	if followerID == followingID {
		return errors.BadRequest("Cannot follow yourself")
	}

	// Check if both users exist and are active
	follower, err := s.userRepo.GetByID(ctx, followerID)
	if err != nil {
		return errors.NotFound("Follower not found")
	}
	if !follower.IsActive {
		return errors.Forbidden("Follower account is inactive")
	}

	following, err := s.userRepo.GetByID(ctx, followingID)
	if err != nil {
		return errors.NotFound("User to follow not found")
	}
	if !following.IsActive {
		return errors.Forbidden("Cannot follow inactive user")
	}

	// Check if already following
	isFollowing, err := s.userRepo.IsFollowing(ctx, followerID, followingID)
	if err != nil {
		return errors.InternalServerError("Failed to check follow status")
	}
	if isFollowing {
		return errors.Conflict("Already following this user")
	}

	return s.userRepo.Follow(ctx, followerID, followingID)
}

func (s *userDomainService) UnfollowUser(ctx context.Context, followerID, followingID uuid.UUID) error {
	// Check if currently following
	isFollowing, err := s.userRepo.IsFollowing(ctx, followerID, followingID)
	if err != nil {
		return errors.InternalServerError("Failed to check follow status")
	}
	if !isFollowing {
		return errors.BadRequest("Not following this user")
	}

	return s.userRepo.Unfollow(ctx, followerID, followingID)
}
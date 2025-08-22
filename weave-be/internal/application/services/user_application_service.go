package services

import (
	"context"

	"github.com/google/uuid"
	"weave-module/errors"
	"weave-be/internal/application/dto"
	"weave-be/internal/domain/entities"
	"weave-be/internal/domain/repositories"
	"weave-be/internal/domain/services"
)

// UserApplicationService handles application-level user operations
// This follows the Application Service pattern and orchestrates domain services
type UserApplicationService interface {
	RegisterUser(ctx context.Context, req dto.RegisterUserRequest) (*dto.UserResponse, error)
	LoginUser(ctx context.Context, req dto.LoginUserRequest) (*dto.LoginResponse, error)
	GetUserProfile(ctx context.Context, userID uuid.UUID) (*dto.UserProfileResponse, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*dto.UserResponse, error)
	UpdateUserProfile(ctx context.Context, userID uuid.UUID, req dto.UpdateUserProfileRequest) (*dto.UserResponse, error)
	FollowUser(ctx context.Context, followerID, followingID uuid.UUID) error
	UnfollowUser(ctx context.Context, followerID, followingID uuid.UUID) error
	GetFollowers(ctx context.Context, userID uuid.UUID, page, limit int) (*dto.PaginatedUsersResponse, error)
	GetFollowing(ctx context.Context, userID uuid.UUID, page, limit int) (*dto.PaginatedUsersResponse, error)
	SearchUsers(ctx context.Context, query string, page, limit int) (*dto.PaginatedUsersResponse, error)
}

type userApplicationService struct {
	userRepo        repositories.UserRepository
	userDomainSvc   services.UserDomainService
}

// NewUserApplicationService creates a new user application service
func NewUserApplicationService(
	userRepo repositories.UserRepository,
	userDomainSvc services.UserDomainService,
) UserApplicationService {
	return &userApplicationService{
		userRepo:      userRepo,
		userDomainSvc: userDomainSvc,
	}
}

func (s *userApplicationService) RegisterUser(ctx context.Context, req dto.RegisterUserRequest) (*dto.UserResponse, error) {
	// Validate request
	if err := req.Validate(); err != nil {
		return nil, errors.BadRequestWithDetails("Invalid registration data", err.Error())
	}

	// Use domain service to handle business logic
	user, err := s.userDomainSvc.RegisterUser(ctx, req.Username, req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	// Convert to DTO
	return dto.UserToResponse(user), nil
}

func (s *userApplicationService) LoginUser(ctx context.Context, req dto.LoginUserRequest) (*dto.LoginResponse, error) {
	// Validate request
	if err := req.Validate(); err != nil {
		return nil, errors.BadRequestWithDetails("Invalid login data", err.Error())
	}

	// Use domain service for authentication
	user, token, err := s.userDomainSvc.AuthenticateUser(ctx, req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	// Convert to DTO
	return &dto.LoginResponse{
		User:  *dto.UserToResponse(user),
		Token: token,
	}, nil
}

func (s *userApplicationService) GetUserProfile(ctx context.Context, userID uuid.UUID) (*dto.UserProfileResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.NotFound("User not found")
	}

	// Get additional profile data
	followersCount, _ := s.userRepo.GetFollowersCount(ctx, userID)
	followingCount, _ := s.userRepo.GetFollowingCount(ctx, userID)
	
	// TODO: Get weaves count, contributions count etc.
	
	return &dto.UserProfileResponse{
		User:           *dto.UserToResponse(user),
		FollowersCount: int(followersCount),
		FollowingCount: int(followingCount),
		WeavesCount:    0, // TODO: implement
		ContributionsCount: 0, // TODO: implement
	}, nil
}

func (s *userApplicationService) GetUserByID(ctx context.Context, userID uuid.UUID) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.NotFound("User not found")
	}

	return dto.UserToResponse(user), nil
}

func (s *userApplicationService) UpdateUserProfile(ctx context.Context, userID uuid.UUID, req dto.UpdateUserProfileRequest) (*dto.UserResponse, error) {
	// Validate request
	if err := req.Validate(); err != nil {
		return nil, errors.BadRequestWithDetails("Invalid profile update data", err.Error())
	}

	// Use domain service
	err := s.userDomainSvc.UpdateUserProfile(ctx, userID, req.ProfileImage, req.Bio)
	if err != nil {
		return nil, err
	}

	// Get updated user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.InternalServerError("Failed to get updated user")
	}

	return dto.UserToResponse(user), nil
}

func (s *userApplicationService) FollowUser(ctx context.Context, followerID, followingID uuid.UUID) error {
	return s.userDomainSvc.FollowUser(ctx, followerID, followingID)
}

func (s *userApplicationService) UnfollowUser(ctx context.Context, followerID, followingID uuid.UUID) error {
	return s.userDomainSvc.UnfollowUser(ctx, followerID, followingID)
}

func (s *userApplicationService) GetFollowers(ctx context.Context, userID uuid.UUID, page, limit int) (*dto.PaginatedUsersResponse, error) {
	offset := (page - 1) * limit
	
	users, err := s.userRepo.GetFollowers(ctx, userID, limit, offset)
	if err != nil {
		return nil, errors.InternalServerError("Failed to get followers")
	}

	total, err := s.userRepo.GetFollowersCount(ctx, userID)
	if err != nil {
		return nil, errors.InternalServerError("Failed to count followers")
	}

	userResponses := make([]dto.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = *dto.UserToResponse(user)
	}

	return &dto.PaginatedUsersResponse{
		Users: userResponses,
		Page:  page,
		Limit: limit,
		Total: int(total),
	}, nil
}

func (s *userApplicationService) GetFollowing(ctx context.Context, userID uuid.UUID, page, limit int) (*dto.PaginatedUsersResponse, error) {
	offset := (page - 1) * limit
	
	users, err := s.userRepo.GetFollowing(ctx, userID, limit, offset)
	if err != nil {
		return nil, errors.InternalServerError("Failed to get following")
	}

	total, err := s.userRepo.GetFollowingCount(ctx, userID)
	if err != nil {
		return nil, errors.InternalServerError("Failed to count following")
	}

	userResponses := make([]dto.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = *dto.UserToResponse(user)
	}

	return &dto.PaginatedUsersResponse{
		Users: userResponses,
		Page:  page,
		Limit: limit,
		Total: int(total),
	}, nil
}

func (s *userApplicationService) SearchUsers(ctx context.Context, query string, page, limit int) (*dto.PaginatedUsersResponse, error) {
	offset := (page - 1) * limit
	
	users, err := s.userRepo.SearchByUsername(ctx, query, limit, offset)
	if err != nil {
		return nil, errors.InternalServerError("Failed to search users")
	}

	// TODO: Get total count for search results
	total := int64(len(users))

	userResponses := make([]dto.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = *dto.UserToResponse(user)
	}

	return &dto.PaginatedUsersResponse{
		Users: userResponses,
		Page:  page,
		Limit: limit,
		Total: int(total),
	}, nil
}
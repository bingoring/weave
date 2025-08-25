package user

import (
	"context"

	"weave-module/errors"
	"weave-be/internal/application/dto"
	"weave-be/internal/application/queries"
	"weave-be/internal/domain/repositories"
)

// SearchUsersUseCase handles user search business logic
type SearchUsersUseCase struct {
	userRepo repositories.UserRepository
}

// NewSearchUsersUseCase creates a new SearchUsersUseCase
func NewSearchUsersUseCase(userRepo repositories.UserRepository) *SearchUsersUseCase {
	return &SearchUsersUseCase{
		userRepo: userRepo,
	}
}

// Execute searches for users by username
func (uc *SearchUsersUseCase) Execute(ctx context.Context, query queries.SearchUsersQuery) (*dto.PaginatedUsersResponse, error) {
	offset := (query.Page - 1) * query.Limit
	
	users, err := uc.userRepo.SearchByUsername(ctx, query.Query, query.Limit, offset)
	if err != nil {
		return nil, errors.InternalServerError("Failed to search users")
	}

	total, err := uc.userRepo.SearchByUsernameCount(ctx, query.Query)
	if err != nil {
		return nil, errors.InternalServerError("Failed to count search results")
	}

	userResponses := make([]dto.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = *dto.UserToResponse(user)
	}

	return &dto.PaginatedUsersResponse{
		Users: userResponses,
		Page:  query.Page,
		Limit: query.Limit,
		Total: int(total),
	}, nil
}

// GetFollowersUseCase handles getting user's followers
type GetFollowersUseCase struct {
	userRepo repositories.UserRepository
}

// NewGetFollowersUseCase creates a new GetFollowersUseCase
func NewGetFollowersUseCase(userRepo repositories.UserRepository) *GetFollowersUseCase {
	return &GetFollowersUseCase{
		userRepo: userRepo,
	}
}

// Execute gets followers for a user
func (uc *GetFollowersUseCase) Execute(ctx context.Context, query queries.GetFollowersQuery) (*dto.PaginatedUsersResponse, error) {
	offset := (query.Page - 1) * query.Limit
	
	users, err := uc.userRepo.GetFollowers(ctx, query.UserID, query.Limit, offset)
	if err != nil {
		return nil, errors.InternalServerError("Failed to get followers")
	}

	total, err := uc.userRepo.GetFollowersCount(ctx, query.UserID)
	if err != nil {
		return nil, errors.InternalServerError("Failed to count followers")
	}

	userResponses := make([]dto.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = *dto.UserToResponse(user)
	}

	return &dto.PaginatedUsersResponse{
		Users: userResponses,
		Page:  query.Page,
		Limit: query.Limit,
		Total: int(total),
	}, nil
}

// GetFollowingUseCase handles getting users that a user is following
type GetFollowingUseCase struct {
	userRepo repositories.UserRepository
}

// NewGetFollowingUseCase creates a new GetFollowingUseCase
func NewGetFollowingUseCase(userRepo repositories.UserRepository) *GetFollowingUseCase {
	return &GetFollowingUseCase{
		userRepo: userRepo,
	}
}

// Execute gets users that a user is following
func (uc *GetFollowingUseCase) Execute(ctx context.Context, query queries.GetFollowingQuery) (*dto.PaginatedUsersResponse, error) {
	offset := (query.Page - 1) * query.Limit
	
	users, err := uc.userRepo.GetFollowing(ctx, query.UserID, query.Limit, offset)
	if err != nil {
		return nil, errors.InternalServerError("Failed to get following")
	}

	total, err := uc.userRepo.GetFollowingCount(ctx, query.UserID)
	if err != nil {
		return nil, errors.InternalServerError("Failed to count following")
	}

	userResponses := make([]dto.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = *dto.UserToResponse(user)
	}

	return &dto.PaginatedUsersResponse{
		Users: userResponses,
		Page:  query.Page,
		Limit: query.Limit,
		Total: int(total),
	}, nil
}
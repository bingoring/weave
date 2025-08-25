package user

import (
	"context"

	"weave-module/errors"
	"weave-be/internal/application/dto"
	"weave-be/internal/application/queries"
	"weave-be/internal/domain/repositories"
)

// GetUserProfileUseCase handles retrieving user profile with additional data
type GetUserProfileUseCase struct {
	userRepo  repositories.UserRepository
	weaveRepo repositories.WeaveRepository
}

// NewGetUserProfileUseCase creates a new GetUserProfileUseCase
func NewGetUserProfileUseCase(userRepo repositories.UserRepository, weaveRepo repositories.WeaveRepository) *GetUserProfileUseCase {
	return &GetUserProfileUseCase{
		userRepo:  userRepo,
		weaveRepo: weaveRepo,
	}
}

// Execute retrieves user profile with additional statistics
func (uc *GetUserProfileUseCase) Execute(ctx context.Context, query queries.GetUserProfileQuery) (*dto.UserProfileResponse, error) {
	user, err := uc.userRepo.GetByID(ctx, query.UserID)
	if err != nil {
		return nil, errors.NotFound("User not found")
	}

	// Get additional profile data
	followersCount, _ := uc.userRepo.GetFollowersCount(ctx, query.UserID)
	followingCount, _ := uc.userRepo.GetFollowingCount(ctx, query.UserID)

	var weavesCount int64 = 0
	if uc.weaveRepo != nil {
		weavesCount, _ = uc.weaveRepo.CountByUser(ctx, query.UserID)
	}

	// Calculate contributions count (forked weaves + collaborative contributions + liked weaves)
	var totalContributions int64 = 0
	if uc.weaveRepo != nil {
		forkedCount, _ := uc.weaveRepo.CountForkedByUser(ctx, query.UserID)
		likedCount, _ := uc.weaveRepo.CountLikedByUser(ctx, query.UserID)
		contributionsCount, _ := uc.weaveRepo.CountContributionsByUser(ctx, query.UserID)
		totalContributions = forkedCount + likedCount + contributionsCount
	}

	return &dto.UserProfileResponse{
		User:               *dto.UserToResponse(user),
		FollowersCount:     int(followersCount),
		FollowingCount:     int(followingCount),
		WeavesCount:        int(weavesCount),
		ContributionsCount: int(totalContributions),
	}, nil
}
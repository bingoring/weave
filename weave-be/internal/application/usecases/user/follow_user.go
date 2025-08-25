package user

import (
	"context"

	"weave-be/internal/application/commands"
	"weave-be/internal/domain/services"
)

// FollowUserUseCase handles user following business logic
type FollowUserUseCase struct {
	userDomainService services.UserDomainService
}

// NewFollowUserUseCase creates a new FollowUserUseCase
func NewFollowUserUseCase(userDomainService services.UserDomainService) *FollowUserUseCase {
	return &FollowUserUseCase{
		userDomainService: userDomainService,
	}
}

// Execute follows a user
func (uc *FollowUserUseCase) Execute(ctx context.Context, cmd commands.FollowUserCommand) error {
	return uc.userDomainService.FollowUser(ctx, cmd.FollowerID, cmd.FollowingID)
}

// UnfollowUserUseCase handles user unfollowing business logic
type UnfollowUserUseCase struct {
	userDomainService services.UserDomainService
}

// NewUnfollowUserUseCase creates a new UnfollowUserUseCase
func NewUnfollowUserUseCase(userDomainService services.UserDomainService) *UnfollowUserUseCase {
	return &UnfollowUserUseCase{
		userDomainService: userDomainService,
	}
}

// Execute unfollows a user
func (uc *UnfollowUserUseCase) Execute(ctx context.Context, cmd commands.UnfollowUserCommand) error {
	return uc.userDomainService.UnfollowUser(ctx, cmd.FollowerID, cmd.FollowingID)
}
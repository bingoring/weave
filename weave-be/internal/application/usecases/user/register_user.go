package user

import (
	"context"

	"weave-be/internal/application/commands"
	"weave-be/internal/application/dto"
	"weave-be/internal/domain/services"
)

// RegisterUserUseCase handles user registration business logic
type RegisterUserUseCase struct {
	userDomainService services.UserDomainService
}

// NewRegisterUserUseCase creates a new RegisterUserUseCase
func NewRegisterUserUseCase(userDomainService services.UserDomainService) *RegisterUserUseCase {
	return &RegisterUserUseCase{
		userDomainService: userDomainService,
	}
}

// Execute registers a new user
func (uc *RegisterUserUseCase) Execute(ctx context.Context, cmd commands.RegisterUserCommand) (*dto.UserResponse, error) {
	// Use domain service to handle business logic
	user, err := uc.userDomainService.RegisterUser(ctx, cmd.Username, cmd.Email, cmd.Password)
	if err != nil {
		return nil, err
	}

	// Convert to DTO
	return dto.UserToResponse(user), nil
}
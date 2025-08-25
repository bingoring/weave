package user

import (
	"context"

	"weave-be/internal/application/commands"
	"weave-be/internal/application/dto"
	"weave-be/internal/domain/services"
)

// LoginUserUseCase handles user authentication business logic
type LoginUserUseCase struct {
	userDomainService services.UserDomainService
}

// NewLoginUserUseCase creates a new LoginUserUseCase
func NewLoginUserUseCase(userDomainService services.UserDomainService) *LoginUserUseCase {
	return &LoginUserUseCase{
		userDomainService: userDomainService,
	}
}

// Execute authenticates a user and returns login response
func (uc *LoginUserUseCase) Execute(ctx context.Context, cmd commands.LoginUserCommand) (*dto.LoginResponse, error) {
	// Use domain service for authentication
	user, token, err := uc.userDomainService.AuthenticateUser(ctx, cmd.Email, cmd.Password)
	if err != nil {
		return nil, err
	}

	// Convert to DTO
	return &dto.LoginResponse{
		User:  *dto.UserToResponse(user),
		Token: token,
	}, nil
}
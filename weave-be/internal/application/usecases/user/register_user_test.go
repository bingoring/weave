package user

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"weave-be/internal/application/commands"
	"weave-be/internal/domain/entities"
)

// Mock UserDomainService for testing
type mockUserDomainService struct {
	registerUserFunc func(ctx context.Context, username, email, password string) (*entities.User, error)
}

func (m *mockUserDomainService) RegisterUser(ctx context.Context, username, email, password string) (*entities.User, error) {
	if m.registerUserFunc != nil {
		return m.registerUserFunc(ctx, username, email, password)
	}
	return nil, errors.New("not implemented")
}

func (m *mockUserDomainService) AuthenticateUser(ctx context.Context, email, password string) (*entities.User, string, error) {
	return nil, "", errors.New("not implemented")
}

func (m *mockUserDomainService) ValidateUserRegistration(ctx context.Context, username, email string) error {
	return errors.New("not implemented")
}

func (m *mockUserDomainService) CanUserCreateWeave(ctx context.Context, userID uuid.UUID) (bool, error) {
	return false, errors.New("not implemented")
}

func (m *mockUserDomainService) UpdateUserProfile(ctx context.Context, userID uuid.UUID, profileImage *string, bio *string) error {
	return errors.New("not implemented")
}

func (m *mockUserDomainService) FollowUser(ctx context.Context, followerID, followingID uuid.UUID) error {
	return errors.New("not implemented")
}

func (m *mockUserDomainService) UnfollowUser(ctx context.Context, followerID, followingID uuid.UUID) error {
	return errors.New("not implemented")
}

func TestRegisterUserUseCase_Execute(t *testing.T) {
	ctx := context.Background()

	t.Run("successful registration", func(t *testing.T) {
		expectedUser := &entities.User{
			ID:       uuid.New(),
			Username: "testuser",
			Email:    "test@example.com",
		}

		mockService := &mockUserDomainService{
			registerUserFunc: func(ctx context.Context, username, email, password string) (*entities.User, error) {
				if username == "testuser" && email == "test@example.com" && password == "password123" {
					return expectedUser, nil
				}
				return nil, errors.New("unexpected parameters")
			},
		}

		useCase := NewRegisterUserUseCase(mockService)
		cmd := commands.RegisterUserCommand{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "password123",
		}

		result, err := useCase.Execute(ctx, cmd)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if result == nil {
			t.Fatal("Expected result, got nil")
		}

		if result.ID != expectedUser.ID {
			t.Errorf("Expected user ID %s, got %s", expectedUser.ID, result.ID)
		}

		if result.Username != expectedUser.Username {
			t.Errorf("Expected username %s, got %s", expectedUser.Username, result.Username)
		}
	})

	t.Run("registration failure", func(t *testing.T) {
		expectedError := errors.New("registration failed")

		mockService := &mockUserDomainService{
			registerUserFunc: func(ctx context.Context, username, email, password string) (*entities.User, error) {
				return nil, expectedError
			},
		}

		useCase := NewRegisterUserUseCase(mockService)
		cmd := commands.RegisterUserCommand{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "password123",
		}

		result, err := useCase.Execute(ctx, cmd)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		if result != nil {
			t.Errorf("Expected nil result, got %v", result)
		}

		if err != expectedError {
			t.Errorf("Expected error %v, got %v", expectedError, err)
		}
	})
}
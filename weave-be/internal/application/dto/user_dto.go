package dto

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"weave-be/internal/domain/entities"
)

// Request DTOs
type RegisterUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

func (r RegisterUserRequest) Validate() error {
	if strings.TrimSpace(r.Username) == "" {
		return fmt.Errorf("username is required")
	}
	if len(r.Username) < 3 || len(r.Username) > 50 {
		return fmt.Errorf("username must be between 3 and 50 characters")
	}
	if strings.TrimSpace(r.Email) == "" {
		return fmt.Errorf("email is required")
	}
	if !isValidEmail(r.Email) {
		return fmt.Errorf("invalid email format")
	}
	if len(r.Password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}
	return nil
}

type LoginUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (r LoginUserRequest) Validate() error {
	if strings.TrimSpace(r.Email) == "" {
		return fmt.Errorf("email is required")
	}
	if !isValidEmail(r.Email) {
		return fmt.Errorf("invalid email format")
	}
	if strings.TrimSpace(r.Password) == "" {
		return fmt.Errorf("password is required")
	}
	return nil
}

type UpdateUserProfileRequest struct {
	ProfileImage *string `json:"profile_image"`
	Bio          *string `json:"bio"`
}

func (r UpdateUserProfileRequest) Validate() error {
	if r.Bio != nil && len(*r.Bio) > 500 {
		return fmt.Errorf("bio cannot exceed 500 characters")
	}
	return nil
}

// Response DTOs
type UserResponse struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	ProfileImage *string   `json:"profile_image"`
	Bio          *string   `json:"bio"`
	IsVerified   bool      `json:"is_verified"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type LoginResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}

type UserProfileResponse struct {
	User               UserResponse `json:"user"`
	FollowersCount     int          `json:"followers_count"`
	FollowingCount     int          `json:"following_count"`
	WeavesCount        int          `json:"weaves_count"`
	ContributionsCount int          `json:"contributions_count"`
}

type PaginatedUsersResponse struct {
	Users []UserResponse `json:"users"`
	Page  int            `json:"page"`
	Limit int            `json:"limit"`
	Total int            `json:"total"`
}

// Conversion functions
func UserToResponse(user *entities.User) *UserResponse {
	return &UserResponse{
		ID:           user.ID,
		Username:     user.Username,
		Email:        user.Email,
		ProfileImage: user.ProfileImage,
		Bio:          user.Bio,
		IsVerified:   user.IsVerified,
		IsActive:     user.IsActive,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}
}

func UsersToResponse(users []*entities.User) []UserResponse {
	responses := make([]UserResponse, len(users))
	for i, user := range users {
		responses[i] = *UserToResponse(user)
	}
	return responses
}

// Helper functions
func isValidEmail(email string) bool {
	// Simple email validation
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}
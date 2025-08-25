package commands

import "github.com/google/uuid"

// RegisterUserCommand represents the command to register a new user
type RegisterUserCommand struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// LoginUserCommand represents the command to authenticate a user
type LoginUserCommand struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// UpdateUserProfileCommand represents the command to update user profile
type UpdateUserProfileCommand struct {
	UserID       uuid.UUID `json:"user_id" validate:"required"`
	ProfileImage *string   `json:"profile_image,omitempty"`
	Bio          *string   `json:"bio,omitempty"`
}

// FollowUserCommand represents the command to follow another user
type FollowUserCommand struct {
	FollowerID  uuid.UUID `json:"follower_id" validate:"required"`
	FollowingID uuid.UUID `json:"following_id" validate:"required"`
}

// UnfollowUserCommand represents the command to unfollow another user
type UnfollowUserCommand struct {
	FollowerID  uuid.UUID `json:"follower_id" validate:"required"`
	FollowingID uuid.UUID `json:"following_id" validate:"required"`
}
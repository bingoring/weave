package queries

import "github.com/google/uuid"

// GetUserByIDQuery represents the query to get user by ID
type GetUserByIDQuery struct {
	UserID uuid.UUID `json:"user_id" validate:"required"`
}

// GetUserProfileQuery represents the query to get user profile with additional data
type GetUserProfileQuery struct {
	UserID uuid.UUID `json:"user_id" validate:"required"`
}

// GetFollowersQuery represents the query to get user's followers
type GetFollowersQuery struct {
	UserID uuid.UUID `json:"user_id" validate:"required"`
	Page   int        `json:"page" validate:"min=1"`
	Limit  int        `json:"limit" validate:"min=1,max=100"`
}

// GetFollowingQuery represents the query to get users that a user is following
type GetFollowingQuery struct {
	UserID uuid.UUID `json:"user_id" validate:"required"`
	Page   int        `json:"page" validate:"min=1"`
	Limit  int        `json:"limit" validate:"min=1,max=100"`
}

// SearchUsersQuery represents the query to search users
type SearchUsersQuery struct {
	Query string `json:"query" validate:"required,min=1"`
	Page  int    `json:"page" validate:"min=1"`
	Limit int    `json:"limit" validate:"min=1,max=100"`
}
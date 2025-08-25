package repositories

import (
	"context"

	"github.com/google/uuid"
	"weave-be/internal/domain/entities"
)

// UserRepository interface defines the contract for user data access
// This follows the Repository Pattern and Dependency Inversion Principle
type UserRepository interface {
	// Create operations
	Create(ctx context.Context, user *entities.User) error
	
	// Read operations
	GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error)
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
	GetByUsername(ctx context.Context, username string) (*entities.User, error)
	List(ctx context.Context, limit, offset int) ([]*entities.User, error)
	Count(ctx context.Context) (int64, error)
	
	// Update operations
	Update(ctx context.Context, user *entities.User) error
	UpdateProfile(ctx context.Context, userID uuid.UUID, profileImage *string, bio *string) error
	UpdateVerificationStatus(ctx context.Context, userID uuid.UUID, isVerified bool) error
	UpdateActiveStatus(ctx context.Context, userID uuid.UUID, isActive bool) error
	
	// Delete operations
	Delete(ctx context.Context, id uuid.UUID) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
	
	// Search operations
	SearchByUsername(ctx context.Context, query string, limit, offset int) ([]*entities.User, error)
	SearchByUsernameCount(ctx context.Context, query string) (int64, error)
	SearchByEmail(ctx context.Context, query string, limit, offset int) ([]*entities.User, error)
	
	// Existence checks
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	ExistsByUsername(ctx context.Context, username string) (bool, error)
	
	// OAuth operations
	GetByGoogleID(ctx context.Context, googleID string) (*entities.User, error)
	ExistsByGoogleID(ctx context.Context, googleID string) (bool, error)
	
	// Follow system
	Follow(ctx context.Context, followerID, followingID uuid.UUID) error
	Unfollow(ctx context.Context, followerID, followingID uuid.UUID) error
	IsFollowing(ctx context.Context, followerID, followingID uuid.UUID) (bool, error)
	GetFollowers(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.User, error)
	GetFollowing(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.User, error)
	GetFollowersCount(ctx context.Context, userID uuid.UUID) (int64, error)
	GetFollowingCount(ctx context.Context, userID uuid.UUID) (int64, error)
}
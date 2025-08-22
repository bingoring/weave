package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"weave-be/internal/domain/entities"
)

// WeaveRepository interface for weave data access operations
type WeaveRepository interface {
	// Create operations
	Create(ctx context.Context, weave *entities.Weave) error
	Fork(ctx context.Context, originalID, newUserID uuid.UUID) (*entities.Weave, error)
	
	// Read operations
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Weave, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.Weave, error)
	GetByChannelID(ctx context.Context, channelID uuid.UUID, limit, offset int) ([]*entities.Weave, error)
	GetPublished(ctx context.Context, limit, offset int) ([]*entities.Weave, error)
	GetFeatured(ctx context.Context, limit, offset int) ([]*entities.Weave, error)
	GetDrafts(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.Weave, error)
	GetForked(ctx context.Context, parentID uuid.UUID, limit, offset int) ([]*entities.Weave, error)
	
	// Update operations
	Update(ctx context.Context, weave *entities.Weave) error
	UpdateContent(ctx context.Context, weaveID uuid.UUID, content entities.WeaveContent) error
	UpdatePublishStatus(ctx context.Context, weaveID uuid.UUID, isPublished bool) error
	UpdateFeaturedStatus(ctx context.Context, weaveID uuid.UUID, isFeatured bool) error
	IncrementViewCount(ctx context.Context, weaveID uuid.UUID) error
	IncrementLikeCount(ctx context.Context, weaveID uuid.UUID) error
	DecrementLikeCount(ctx context.Context, weaveID uuid.UUID) error
	IncrementForkCount(ctx context.Context, weaveID uuid.UUID) error
	
	// Delete operations
	Delete(ctx context.Context, id uuid.UUID) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
	
	// Search operations
	Search(ctx context.Context, query string, channelID *uuid.UUID, limit, offset int) ([]*entities.Weave, error)
	SearchByTags(ctx context.Context, tags []string, limit, offset int) ([]*entities.Weave, error)
	
	// Analytics
	Count(ctx context.Context) (int64, error)
	CountByChannel(ctx context.Context, channelID uuid.UUID) (int64, error)
	CountByUser(ctx context.Context, userID uuid.UUID) (int64, error)
	GetTrending(ctx context.Context, timeframe string, limit, offset int) ([]*entities.Weave, error)
	GetPopular(ctx context.Context, limit, offset int) ([]*entities.Weave, error)
	
	// Like system
	Like(ctx context.Context, weaveID, userID uuid.UUID) error
	Unlike(ctx context.Context, weaveID, userID uuid.UUID) error
	IsLiked(ctx context.Context, weaveID, userID uuid.UUID) (bool, error)
	GetLikedBy(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.Weave, error)
	
	// Version control
	CreateVersion(ctx context.Context, weaveID uuid.UUID, content entities.WeaveContent, changeLog *string) error
	GetVersions(ctx context.Context, weaveID uuid.UUID) ([]*entities.WeaveVersion, error)
	GetVersion(ctx context.Context, weaveID uuid.UUID, version int) (*entities.WeaveVersion, error)
}

// WeaveVersion represents a historical version of a weave
type WeaveVersion struct {
	ID        uuid.UUID
	WeaveID   uuid.UUID
	Version   int
	Title     string
	Content   entities.WeaveContent
	ChangeLog *string
	CreatedAt time.Time
}
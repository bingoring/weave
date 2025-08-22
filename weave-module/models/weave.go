package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WeaveStatus string

const (
	WeaveStatusDraft       WeaveStatus = "draft"
	WeaveStatusInReview    WeaveStatus = "in_review"
	WeaveStatusPublished   WeaveStatus = "published"
	WeaveStatusArchived    WeaveStatus = "archived"
	WeaveStatusDeleted     WeaveStatus = "deleted"
)

type WeaveType string

const (
	WeaveTypeOriginal     WeaveType = "original"
	WeaveTypeFork         WeaveType = "fork"
	WeaveTypeContribution WeaveType = "contribution"
	WeaveTypeMerge        WeaveType = "merge"
)

type Weave struct {
	ID                  uuid.UUID   `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID              uuid.UUID   `gorm:"type:uuid;not null;index" json:"user_id"`
	ChannelID           uuid.UUID   `gorm:"type:uuid;not null;index" json:"channel_id"`
	Title               string      `gorm:"not null;size:200" json:"title"`
	Description         *string     `gorm:"type:text" json:"description"`
	CoverImage          *string     `gorm:"size:500" json:"cover_image"`
	Content             string      `gorm:"type:jsonb;not null" json:"content"`
	Status              WeaveStatus `gorm:"type:varchar(20);default:'draft';index" json:"status"`
	Type                WeaveType   `gorm:"type:varchar(20);default:'original';index" json:"type"`
	Version             int         `gorm:"default:1" json:"version"`
	ParentWeaveID       *uuid.UUID  `gorm:"type:uuid;index" json:"parent_weave_id"`
	OriginalWeaveID     *uuid.UUID  `gorm:"type:uuid;index" json:"original_weave_id"`
	IsCollaborationOpen bool        `gorm:"default:true" json:"is_collaboration_open"`
	IsFeatured          bool        `gorm:"default:false;index" json:"is_featured"`
	ViewCount           int         `gorm:"default:0" json:"view_count"`
	LikeCount           int         `gorm:"default:0" json:"like_count"`
	ForkCount           int         `gorm:"default:0" json:"fork_count"`
	ContributionCount   int         `gorm:"default:0" json:"contribution_count"`
	PublishedAt         *time.Time  `json:"published_at"`
	CreatedAt           time.Time   `gorm:"autoCreateTime;index" json:"created_at"`
	UpdatedAt           time.Time   `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	User           User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Channel        Channel        `gorm:"foreignKey:ChannelID" json:"channel,omitempty"`
	ParentWeave    *Weave         `gorm:"foreignKey:ParentWeaveID" json:"parent_weave,omitempty"`
	OriginalWeave  *Weave         `gorm:"foreignKey:OriginalWeaveID" json:"original_weave,omitempty"`
	ChildWeaves    []Weave        `gorm:"foreignKey:ParentWeaveID" json:"child_weaves,omitempty"`
	Forks          []Weave        `gorm:"foreignKey:OriginalWeaveID" json:"forks,omitempty"`
	Likes          []WeaveLike    `gorm:"foreignKey:WeaveID" json:"likes,omitempty"`
	Versions       []WeaveVersion `gorm:"foreignKey:WeaveID" json:"versions,omitempty"`
	LabComments    []LabComment   `gorm:"foreignKey:WeaveID" json:"lab_comments,omitempty"`
	Contributions  []Contribution `gorm:"foreignKey:WeaveID" json:"contributions,omitempty"`
	Tags           []WeaveTag     `gorm:"many2many:weave_tag_relations;" json:"tags,omitempty"`
	Collections    []WeaveCollection `gorm:"many2many:collection_weaves;" json:"collections,omitempty"`
}

type WeaveVersion struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	WeaveID     uuid.UUID `gorm:"type:uuid;not null;index" json:"weave_id"`
	UserID      uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	Version     int       `gorm:"not null;index" json:"version"`
	Title       string    `gorm:"not null;size:200" json:"title"`
	Description *string   `gorm:"type:text" json:"description"`
	Content     string    `gorm:"type:jsonb;not null" json:"content"`
	ChangeLog   *string   `gorm:"type:text" json:"change_log"`
	ContentDiff *string   `gorm:"type:jsonb" json:"content_diff"`
	IsMajor     bool      `gorm:"default:false" json:"is_major"`
	CreatedAt   time.Time `gorm:"autoCreateTime;index" json:"created_at"`

	// Relationships
	Weave Weave `gorm:"foreignKey:WeaveID" json:"weave,omitempty"`
	User  User  `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

type WeaveTimeline struct {
	ID          uuid.UUID         `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	WeaveID     uuid.UUID         `gorm:"type:uuid;not null;index" json:"weave_id"`
	UserID      uuid.UUID         `gorm:"type:uuid;not null;index" json:"user_id"`
	EventType   WeaveTimelineType `gorm:"type:varchar(30);not null;index" json:"event_type"`
	Title       string            `gorm:"not null;size:200" json:"title"`
	Description *string           `gorm:"type:text" json:"description"`
	Metadata    *string           `gorm:"type:jsonb" json:"metadata"`
	CreatedAt   time.Time         `gorm:"autoCreateTime;index" json:"created_at"`

	// Relationships
	Weave Weave `gorm:"foreignKey:WeaveID" json:"weave,omitempty"`
	User  User  `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

type WeaveTimelineType string

const (
	TimelineCreated           WeaveTimelineType = "created"
	TimelineUpdated           WeaveTimelineType = "updated"
	TimelinePublished         WeaveTimelineType = "published"
	TimelineForked            WeaveTimelineType = "forked"
	TimelineContributionAdded WeaveTimelineType = "contribution_added"
	TimelineContributionMerged WeaveTimelineType = "contribution_merged"
	TimelineStatusChanged     WeaveTimelineType = "status_changed"
	TimelineCommentAdded      WeaveTimelineType = "comment_added"
	TimelineLiked             WeaveTimelineType = "liked"
	TimelineCollectionAdded   WeaveTimelineType = "collection_added"
)

type WeaveLike struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	WeaveID   uuid.UUID `gorm:"type:uuid;not null" json:"weave_id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`

	// Relationships
	User  User  `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Weave Weave `gorm:"foreignKey:WeaveID" json:"weave,omitempty"`
}

type WeaveTag struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name        string    `gorm:"uniqueIndex;not null;size:50" json:"name"`
	Description *string   `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`

	// Relationships
	Weaves []Weave `gorm:"many2many:weave_tag_relations;" json:"weaves,omitempty"`
}

type WeaveCollection struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID      uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Name        string    `gorm:"not null;size:100" json:"name"`
	Description *string   `gorm:"type:text" json:"description"`
	IsPublic    bool      `gorm:"default:false" json:"is_public"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	User   User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Weaves []Weave `gorm:"many2many:collection_weaves;" json:"weaves,omitempty"`
}

func (w *Weave) BeforeCreate(tx *gorm.DB) error {
	if w.ID == uuid.Nil {
		w.ID = uuid.New()
	}
	return nil
}

func (wv *WeaveVersion) BeforeCreate(tx *gorm.DB) error {
	if wv.ID == uuid.Nil {
		wv.ID = uuid.New()
	}
	return nil
}

func (wt *WeaveTimeline) BeforeCreate(tx *gorm.DB) error {
	if wt.ID == uuid.Nil {
		wt.ID = uuid.New()
	}
	return nil
}

func (wl *WeaveLike) BeforeCreate(tx *gorm.DB) error {
	if wl.ID == uuid.Nil {
		wl.ID = uuid.New()
	}
	return nil
}

func (wt *WeaveTag) BeforeCreate(tx *gorm.DB) error {
	if wt.ID == uuid.Nil {
		wt.ID = uuid.New()
	}
	return nil
}

func (wc *WeaveCollection) BeforeCreate(tx *gorm.DB) error {
	if wc.ID == uuid.Nil {
		wc.ID = uuid.New()
	}
	return nil
}
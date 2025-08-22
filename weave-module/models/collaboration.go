package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ContributionType string

const (
	ContributionTypeSuggestion   ContributionType = "suggestion"
	ContributionTypeContentEdit  ContributionType = "content_edit"
	ContributionTypeStructural   ContributionType = "structural"
	ContributionTypeFork         ContributionType = "fork"
	ContributionTypeMergeRequest ContributionType = "merge_request"
)

type ContributionStatus string

const (
	ContributionStatusPending  ContributionStatus = "pending"
	ContributionStatusReviewing ContributionStatus = "reviewing"
	ContributionStatusAccepted ContributionStatus = "accepted"
	ContributionStatusRejected ContributionStatus = "rejected"
	ContributionStatusMerged   ContributionStatus = "merged"
)

type Contribution struct {
	ID               uuid.UUID          `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID           uuid.UUID          `gorm:"type:uuid;not null;index" json:"user_id"`
	WeaveID          uuid.UUID          `gorm:"type:uuid;not null;index" json:"weave_id"`
	Type             ContributionType   `gorm:"type:varchar(30);not null;index" json:"type"`
	Title            string             `gorm:"not null;size:200" json:"title"`
	Description      *string            `gorm:"type:text" json:"description"`
	OriginalContent  *string            `gorm:"type:jsonb" json:"original_content"`
	ProposedContent  *string            `gorm:"type:jsonb" json:"proposed_content"`
	ContentDiff      *string            `gorm:"type:jsonb" json:"content_diff"`
	Status           ContributionStatus `gorm:"type:varchar(20);default:'pending';index" json:"status"`
	ReviewerID       *uuid.UUID         `gorm:"type:uuid;index" json:"reviewer_id"`
	ReviewedAt       *time.Time         `json:"reviewed_at"`
	ReviewComment    *string            `gorm:"type:text" json:"review_comment"`
	VoteScore        int                `gorm:"default:0" json:"vote_score"`
	Priority         int                `gorm:"default:0;index" json:"priority"`
	CreatedAt        time.Time          `gorm:"autoCreateTime;index" json:"created_at"`
	UpdatedAt        time.Time          `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	User     User  `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Weave    Weave `gorm:"foreignKey:WeaveID" json:"weave,omitempty"`
	Reviewer *User `gorm:"foreignKey:ReviewerID" json:"reviewer,omitempty"`
	Votes    []ContributionVote `gorm:"foreignKey:ContributionID" json:"votes,omitempty"`
	Comments []ContributionComment `gorm:"foreignKey:ContributionID" json:"comments,omitempty"`
}

type CommentType string

const (
	CommentTypeGeneral    CommentType = "general"
	CommentTypeSuggestion CommentType = "suggestion"
	CommentTypeQuestion   CommentType = "question"
	CommentTypeIssue      CommentType = "issue"
	CommentTypeApproval   CommentType = "approval"
)

type LabComment struct {
	ID               uuid.UUID   `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID           uuid.UUID   `gorm:"type:uuid;not null;index" json:"user_id"`
	WeaveID          uuid.UUID   `gorm:"type:uuid;not null;index" json:"weave_id"`
	ParentCommentID  *uuid.UUID  `gorm:"type:uuid;index" json:"parent_comment_id"`
	Type             CommentType `gorm:"type:varchar(20);default:'general'" json:"type"`
	Content          string      `gorm:"type:text;not null" json:"content"`
	ContentPosition  *string     `gorm:"type:jsonb" json:"content_position"`
	IsResolved       bool        `gorm:"default:false;index" json:"is_resolved"`
	ResolvedBy       *uuid.UUID  `gorm:"type:uuid" json:"resolved_by"`
	ResolvedAt       *time.Time  `json:"resolved_at"`
	LikeCount        int         `gorm:"default:0" json:"like_count"`
	CreatedAt        time.Time   `gorm:"autoCreateTime;index" json:"created_at"`
	UpdatedAt        time.Time   `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	User          User         `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Weave         Weave        `gorm:"foreignKey:WeaveID" json:"weave,omitempty"`
	Resolver      *User        `gorm:"foreignKey:ResolvedBy" json:"resolver,omitempty"`
	ParentComment *LabComment  `gorm:"foreignKey:ParentCommentID" json:"parent_comment,omitempty"`
	Replies       []LabComment `gorm:"foreignKey:ParentCommentID" json:"replies,omitempty"`
}

type ContributionComment struct {
	ID             uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID         uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	ContributionID uuid.UUID `gorm:"type:uuid;not null;index" json:"contribution_id"`
	Content        string    `gorm:"type:text;not null" json:"content"`
	CreatedAt      time.Time `gorm:"autoCreateTime;index" json:"created_at"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	User         User         `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Contribution Contribution `gorm:"foreignKey:ContributionID" json:"contribution,omitempty"`
}

type ContributionVote struct {
	ID             uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID         uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	ContributionID uuid.UUID `gorm:"type:uuid;not null" json:"contribution_id"`
	VoteType       string    `gorm:"not null;size:10" json:"vote_type"` // up, down
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`

	// Relationships
	User         User         `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Contribution Contribution `gorm:"foreignKey:ContributionID" json:"contribution,omitempty"`
}

func (c *Contribution) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

func (lc *LabComment) BeforeCreate(tx *gorm.DB) error {
	if lc.ID == uuid.Nil {
		lc.ID = uuid.New()
	}
	return nil
}

func (cc *ContributionComment) BeforeCreate(tx *gorm.DB) error {
	if cc.ID == uuid.Nil {
		cc.ID = uuid.New()
	}
	return nil
}

func (cv *ContributionVote) BeforeCreate(tx *gorm.DB) error {
	if cv.ID == uuid.Nil {
		cv.ID = uuid.New()
	}
	return nil
}
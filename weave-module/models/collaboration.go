package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Contribution struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID      uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	WeaveID     uuid.UUID `gorm:"type:uuid;not null" json:"weave_id"`
	Type        string    `gorm:"not null;size:50" json:"type"` // suggestion, fork, merge, etc.
	Title       string    `gorm:"not null;size:200" json:"title"`
	Description *string   `gorm:"type:text" json:"description"`
	Content     *string   `gorm:"type:jsonb" json:"content"`
	Status      string    `gorm:"not null;size:20;default:'pending'" json:"status"` // pending, accepted, rejected
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	User  User  `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Weave Weave `gorm:"foreignKey:WeaveID" json:"weave,omitempty"`
}

type LabComment struct {
	ID               uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID           uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	WeaveID          uuid.UUID `gorm:"type:uuid;not null" json:"weave_id"`
	ParentCommentID  *uuid.UUID `gorm:"type:uuid" json:"parent_comment_id"`
	Content          string    `gorm:"type:text;not null" json:"content"`
	IsResolved       bool      `gorm:"default:false" json:"is_resolved"`
	LikeCount        int       `gorm:"default:0" json:"like_count"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	User          User         `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Weave         Weave        `gorm:"foreignKey:WeaveID" json:"weave,omitempty"`
	ParentComment *LabComment  `gorm:"foreignKey:ParentCommentID" json:"parent_comment,omitempty"`
	Replies       []LabComment `gorm:"foreignKey:ParentCommentID" json:"replies,omitempty"`
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

func (cv *ContributionVote) BeforeCreate(tx *gorm.DB) error {
	if cv.ID == uuid.Nil {
		cv.ID = uuid.New()
	}
	return nil
}
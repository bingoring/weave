package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Weave struct {
	ID             uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID         uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	ChannelID      uuid.UUID `gorm:"type:uuid;not null" json:"channel_id"`
	Title          string    `gorm:"not null;size:200" json:"title"`
	CoverImage     *string   `gorm:"size:500" json:"cover_image"`
	Content        string    `gorm:"type:jsonb;not null" json:"content"`
	Version        int       `gorm:"default:1" json:"version"`
	ParentWeaveID  *uuid.UUID `gorm:"type:uuid" json:"parent_weave_id"`
	IsPublished    bool      `gorm:"default:false" json:"is_published"`
	IsFeatured     bool      `gorm:"default:false" json:"is_featured"`
	ViewCount      int       `gorm:"default:0" json:"view_count"`
	LikeCount      int       `gorm:"default:0" json:"like_count"`
	ForkCount      int       `gorm:"default:0" json:"fork_count"`
	CommentCount   int       `gorm:"default:0" json:"comment_count"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	User           User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Channel        Channel        `gorm:"foreignKey:ChannelID" json:"channel,omitempty"`
	ParentWeave    *Weave         `gorm:"foreignKey:ParentWeaveID" json:"parent_weave,omitempty"`
	ChildWeaves    []Weave        `gorm:"foreignKey:ParentWeaveID" json:"child_weaves,omitempty"`
	Likes          []WeaveLike    `gorm:"foreignKey:WeaveID" json:"likes,omitempty"`
	Versions       []WeaveVersion `gorm:"foreignKey:WeaveID" json:"versions,omitempty"`
	LabComments    []LabComment   `gorm:"foreignKey:WeaveID" json:"lab_comments,omitempty"`
	Contributions  []Contribution `gorm:"foreignKey:WeaveID" json:"contributions,omitempty"`
	Tags           []WeaveTag     `gorm:"many2many:weave_tag_relations;" json:"tags,omitempty"`
}

type WeaveVersion struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	WeaveID   uuid.UUID `gorm:"type:uuid;not null" json:"weave_id"`
	Version   int       `gorm:"not null" json:"version"`
	Title     string    `gorm:"not null;size:200" json:"title"`
	Content   string    `gorm:"type:jsonb;not null" json:"content"`
	ChangeLog *string   `gorm:"type:text" json:"change_log"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`

	// Relationships
	Weave Weave `gorm:"foreignKey:WeaveID" json:"weave,omitempty"`
}

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
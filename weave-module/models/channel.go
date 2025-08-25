package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Channel struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name        string    `gorm:"uniqueIndex;not null;size:100" json:"name"`
	Slug        string    `gorm:"uniqueIndex;not null;size:100" json:"slug"`
	Description *string   `gorm:"type:text" json:"description"`
	CoverImage  *string   `gorm:"size:500" json:"cover_image"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	IsPublic    bool      `gorm:"default:true" json:"is_public"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	Weaves []Weave `gorm:"foreignKey:ChannelID" json:"weaves,omitempty"`
}

func (c *Channel) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}
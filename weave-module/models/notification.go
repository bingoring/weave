package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Notification struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Type      string    `gorm:"not null;size:50" json:"type"` // like, comment, follow, contribution, etc.
	Title     string    `gorm:"not null;size:200" json:"title"`
	Message   string    `gorm:"type:text;not null" json:"message"`
	Data      *string   `gorm:"type:jsonb" json:"data"` // Additional data like weave_id, user_id, etc.
	IsRead    bool      `gorm:"default:false" json:"is_read"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`

	// Relationships
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

type NotificationSetting struct {
	ID               uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID           uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	EmailLikes       bool      `gorm:"default:true" json:"email_likes"`
	EmailComments    bool      `gorm:"default:true" json:"email_comments"`
	EmailFollows     bool      `gorm:"default:true" json:"email_follows"`
	EmailContributions bool    `gorm:"default:true" json:"email_contributions"`
	PushLikes        bool      `gorm:"default:true" json:"push_likes"`
	PushComments     bool      `gorm:"default:true" json:"push_comments"`
	PushFollows      bool      `gorm:"default:true" json:"push_follows"`
	PushContributions bool     `gorm:"default:true" json:"push_contributions"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (n *Notification) BeforeCreate(tx *gorm.DB) error {
	if n.ID == uuid.Nil {
		n.ID = uuid.New()
	}
	return nil
}

func (ns *NotificationSetting) BeforeCreate(tx *gorm.DB) error {
	if ns.ID == uuid.Nil {
		ns.ID = uuid.New()
	}
	return nil
}
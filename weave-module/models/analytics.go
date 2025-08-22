package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// WeaveView tracks individual view events for analytics
type WeaveView struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	WeaveID     uuid.UUID `gorm:"type:uuid;not null;index" json:"weave_id"`
	UserID      *uuid.UUID `gorm:"type:uuid;index" json:"user_id"` // nullable for anonymous views
	IPAddress   string    `gorm:"size:45;index" json:"ip_address"`
	UserAgent   *string   `gorm:"type:text" json:"user_agent"`
	Duration    *int      `json:"duration"` // viewing duration in seconds
	ViewedAt    time.Time `gorm:"autoCreateTime;index" json:"viewed_at"`

	// Relationships
	Weave Weave `gorm:"foreignKey:WeaveID" json:"weave,omitempty"`
	User  *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// WeaveAnalytics stores aggregated analytics data for efficient querying
type WeaveAnalytics struct {
	ID                uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	WeaveID           uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"weave_id"`
	TotalViews        int       `gorm:"default:0" json:"total_views"`
	UniqueViews       int       `gorm:"default:0" json:"unique_views"`
	TotalLikes        int       `gorm:"default:0" json:"total_likes"`
	TotalForks        int       `gorm:"default:0" json:"total_forks"`
	TotalContributions int      `gorm:"default:0" json:"total_contributions"`
	TotalComments     int       `gorm:"default:0" json:"total_comments"`
	AvgViewDuration   float64   `gorm:"default:0" json:"avg_view_duration"`
	TrendingScore     float64   `gorm:"default:0;index" json:"trending_score"`
	UpdatedAt         time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	Weave Weave `gorm:"foreignKey:WeaveID" json:"weave,omitempty"`
}

// UserAnalytics stores aggregated user analytics
type UserAnalytics struct {
	ID                    uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID                uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`
	TotalWeaves           int       `gorm:"default:0" json:"total_weaves"`
	TotalPublishedWeaves  int       `gorm:"default:0" json:"total_published_weaves"`
	TotalContributions    int       `gorm:"default:0" json:"total_contributions"`
	TotalLikesReceived    int       `gorm:"default:0" json:"total_likes_received"`
	TotalViewsReceived    int       `gorm:"default:0" json:"total_views_received"`
	TotalFollowers        int       `gorm:"default:0" json:"total_followers"`
	TotalFollowing        int       `gorm:"default:0" json:"total_following"`
	InfluenceScore        float64   `gorm:"default:0;index" json:"influence_score"`
	ContributionScore     float64   `gorm:"default:0;index" json:"contribution_score"`
	LastActiveAt          *time.Time `json:"last_active_at"`
	UpdatedAt             time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// DailyStats stores daily aggregated statistics for time-series analysis
type DailyStats struct {
	ID               uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Date             time.Time `gorm:"type:date;index" json:"date"`
	TotalWeaves      int       `gorm:"default:0" json:"total_weaves"`
	TotalUsers       int       `gorm:"default:0" json:"total_users"`
	TotalViews       int       `gorm:"default:0" json:"total_views"`
	TotalLikes       int       `gorm:"default:0" json:"total_likes"`
	TotalContributions int     `gorm:"default:0" json:"total_contributions"`
	NewUsers         int       `gorm:"default:0" json:"new_users"`
	ActiveUsers      int       `gorm:"default:0" json:"active_users"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// TrendingWeave stores trending weave calculations
type TrendingWeave struct {
	ID            uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	WeaveID       uuid.UUID `gorm:"type:uuid;not null;index" json:"weave_id"`
	Score         float64   `gorm:"not null;index" json:"score"`
	Rank          int       `gorm:"not null;index" json:"rank"`
	Period        string    `gorm:"type:varchar(20);not null;index" json:"period"` // hourly, daily, weekly, monthly
	CalculatedAt  time.Time `gorm:"autoCreateTime;index" json:"calculated_at"`

	// Relationships
	Weave Weave `gorm:"foreignKey:WeaveID" json:"weave,omitempty"`
}

func (wv *WeaveView) BeforeCreate(tx *gorm.DB) error {
	if wv.ID == uuid.Nil {
		wv.ID = uuid.New()
	}
	return nil
}

func (wa *WeaveAnalytics) BeforeCreate(tx *gorm.DB) error {
	if wa.ID == uuid.Nil {
		wa.ID = uuid.New()
	}
	return nil
}

func (ua *UserAnalytics) BeforeCreate(tx *gorm.DB) error {
	if ua.ID == uuid.Nil {
		ua.ID = uuid.New()
	}
	return nil
}

func (ds *DailyStats) BeforeCreate(tx *gorm.DB) error {
	if ds.ID == uuid.Nil {
		ds.ID = uuid.New()
	}
	return nil
}

func (tw *TrendingWeave) BeforeCreate(tx *gorm.DB) error {
	if tw.ID == uuid.Nil {
		tw.ID = uuid.New()
	}
	return nil
}
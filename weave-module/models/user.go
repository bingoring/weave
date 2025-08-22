package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Username     string    `gorm:"uniqueIndex;not null;size:50" json:"username"`
	Email        string    `gorm:"uniqueIndex;not null;size:255" json:"email"`
	PasswordHash string    `gorm:"not null;size:255" json:"-"`
	ProfileImage *string   `gorm:"size:500" json:"profile_image"`
	Bio          *string   `gorm:"type:text" json:"bio"`
	IsVerified   bool      `gorm:"default:false" json:"is_verified"`
	IsActive     bool      `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	Weaves         []Weave         `gorm:"foreignKey:UserID" json:"weaves,omitempty"`
	Likes          []WeaveLike     `gorm:"foreignKey:UserID" json:"likes,omitempty"`
	Follows        []UserFollow    `gorm:"foreignKey:FollowerID" json:"follows,omitempty"`
	Followers      []UserFollow    `gorm:"foreignKey:FollowingID" json:"followers,omitempty"`
	Contributions  []Contribution  `gorm:"foreignKey:UserID" json:"contributions,omitempty"`
	LabComments    []LabComment    `gorm:"foreignKey:UserID" json:"lab_comments,omitempty"`
	Notifications  []Notification  `gorm:"foreignKey:UserID" json:"notifications,omitempty"`
}

type UserFollow struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	FollowerID  uuid.UUID `gorm:"type:uuid;not null" json:"follower_id"`
	FollowingID uuid.UUID `gorm:"type:uuid;not null" json:"following_id"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`

	// Relationships
	Follower  User `gorm:"foreignKey:FollowerID" json:"follower,omitempty"`
	Following User `gorm:"foreignKey:FollowingID" json:"following,omitempty"`
}

type UserProfile struct {
	UserID              uuid.UUID `gorm:"type:uuid;primary_key" json:"user_id"`
	FollowersCount      int       `gorm:"default:0" json:"followers_count"`
	FollowingCount      int       `gorm:"default:0" json:"following_count"`
	WeavesCount         int       `gorm:"default:0" json:"weaves_count"`
	ContributionsCount  int       `gorm:"default:0" json:"contributions_count"`
	TotalLikesReceived  int       `gorm:"default:0" json:"total_likes_received"`
	FeaturedWeavesCount int       `gorm:"default:0" json:"featured_weaves_count"`

	// Relationships
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

func (uf *UserFollow) BeforeCreate(tx *gorm.DB) error {
	if uf.ID == uuid.Nil {
		uf.ID = uuid.New()
	}
	return nil
}
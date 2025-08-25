package entities

import (
	"time"

	"github.com/google/uuid"
)

// User domain entity - represents core business logic
type User struct {
	ID           uuid.UUID
	Username     string
	Email        string
	PasswordHash string
	ProfileImage *string
	Bio          *string
	IsVerified   bool
	IsActive     bool
	// OAuth fields
	GoogleID     *string
	GoogleEmail  *string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// User business methods
func (u *User) IsValidForRegistration() bool {
	return u.Username != "" && u.Email != "" && u.PasswordHash != ""
}

func (u *User) CanCreateWeave() bool {
	return u.IsActive
}

func (u *User) CanModerateContent() bool {
	return u.IsActive && u.IsVerified
	// In the future, add role-based permissions
}

func (u *User) GetDisplayName() string {
	if u.Username != "" {
		return u.Username
	}
	return u.Email
}

func (u *User) Activate() {
	u.IsActive = true
	u.UpdatedAt = time.Now()
}

func (u *User) Deactivate() {
	u.IsActive = false
	u.UpdatedAt = time.Now()
}

func (u *User) Verify() {
	u.IsVerified = true
	u.UpdatedAt = time.Now()
}

// OAuth business methods
func (u *User) IsOAuthUser() bool {
	return u.GoogleID != nil
}

func (u *User) LinkGoogleAccount(googleID, googleEmail string) {
	u.GoogleID = &googleID
	u.GoogleEmail = &googleEmail
	u.IsVerified = true
	u.UpdatedAt = time.Now()
}

func (u *User) UnlinkGoogleAccount() {
	u.GoogleID = nil
	u.GoogleEmail = nil
	u.UpdatedAt = time.Now()
}

func NewUser(username, email, passwordHash string) *User {
	return &User{
		ID:           uuid.New(),
		Username:     username,
		Email:        email,
		PasswordHash: passwordHash,
		IsVerified:   false,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

func NewOAuthUser(username, email, googleID, googleEmail string) *User {
	return &User{
		ID:          uuid.New(),
		Username:    username,
		Email:       email,
		GoogleID:    &googleID,
		GoogleEmail: &googleEmail,
		IsVerified:  true, // OAuth users are auto-verified
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// NewEmailUser creates a new user via email verification (no password)
func NewEmailUser(username, email string) *User {
	return &User{
		ID:        uuid.New(),
		Username:  username,
		Email:     email,
		IsVerified: true, // Email verified users are auto-verified
		IsActive:   true,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

// IsEmailAuth checks if user was created via email authentication
func (u *User) IsEmailAuth() bool {
	return u.PasswordHash == "" && u.GoogleID == nil
}
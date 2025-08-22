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
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// User business methods
func (u *User) IsValidForRegistration() bool {
	return u.Username != "" && u.Email != "" && u.PasswordHash != ""
}

func (u *User) CanCreateWeave() bool {
	return u.IsActive && u.IsVerified
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
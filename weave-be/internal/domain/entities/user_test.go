package entities

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewUser(t *testing.T) {
	username := "testuser"
	email := "test@example.com"
	passwordHash := "hashed_password"

	user := NewUser(username, email, passwordHash)

	if user.Username != username {
		t.Errorf("Expected username %s, got %s", username, user.Username)
	}

	if user.Email != email {
		t.Errorf("Expected email %s, got %s", email, user.Email)
	}

	if user.PasswordHash != passwordHash {
		t.Errorf("Expected password hash %s, got %s", passwordHash, user.PasswordHash)
	}

	if user.IsActive != true {
		t.Error("Expected new user to be active")
	}

	if user.IsVerified != false {
		t.Error("Expected new user to not be verified")
	}

	if user.ID == uuid.Nil {
		t.Error("Expected user ID to be generated")
	}

	if user.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be set")
	}

	if user.UpdatedAt.IsZero() {
		t.Error("Expected UpdatedAt to be set")
	}
}

func TestUser_CanCreateWeave(t *testing.T) {
	user := &User{
		ID:        uuid.New(),
		Username:  "testuser",
		Email:     "test@example.com",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if !user.CanCreateWeave() {
		t.Error("Active user should be able to create weave")
	}

	user.IsActive = false
	if user.CanCreateWeave() {
		t.Error("Inactive user should not be able to create weave")
	}
}

func TestUser_Activate(t *testing.T) {
	user := &User{
		ID:        uuid.New(),
		Username:  "testuser",
		Email:     "test@example.com",
		IsActive:  false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	oldUpdatedAt := user.UpdatedAt
	time.Sleep(time.Millisecond) // Ensure time difference
	
	user.Activate()

	if !user.IsActive {
		t.Error("User should be active after calling Activate")
	}

	if !user.UpdatedAt.After(oldUpdatedAt) {
		t.Error("UpdatedAt should be updated when activating user")
	}
}

func TestUser_Deactivate(t *testing.T) {
	user := &User{
		ID:        uuid.New(),
		Username:  "testuser",
		Email:     "test@example.com",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	oldUpdatedAt := user.UpdatedAt
	time.Sleep(time.Millisecond) // Ensure time difference
	
	user.Deactivate()

	if user.IsActive {
		t.Error("User should be inactive after calling Deactivate")
	}

	if !user.UpdatedAt.After(oldUpdatedAt) {
		t.Error("UpdatedAt should be updated when deactivating user")
	}
}

func TestUser_Verify(t *testing.T) {
	user := &User{
		ID:         uuid.New(),
		Username:   "testuser",
		Email:      "test@example.com",
		IsVerified: false,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	oldUpdatedAt := user.UpdatedAt
	time.Sleep(time.Millisecond) // Ensure time difference
	
	user.Verify()

	if !user.IsVerified {
		t.Error("User should be verified after calling Verify")
	}

	if !user.UpdatedAt.After(oldUpdatedAt) {
		t.Error("UpdatedAt should be updated when verifying user")
	}
}
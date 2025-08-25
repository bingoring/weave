package oauth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"weave-module/config"
)

// OAuthService OAuth service manager
type OAuthService struct {
	factory   *ProviderFactory
	config    config.OAuthConfig
	stateKeys map[string]StateInfo // state token storage (Redis recommended in production)
}

// StateInfo state token information
type StateInfo struct {
	UserID     *uuid.UUID `json:"user_id,omitempty"`
	Provider   string     `json:"provider"`
	Action     string     `json:"action"` // "login" or "connect"
	CreatedAt  time.Time  `json:"created_at"`
	ExpiresAt  time.Time  `json:"expires_at"`
}

// NewOAuthService creates OAuth service
func NewOAuthService(cfg config.OAuthConfig) *OAuthService {
	service := &OAuthService{
		factory:   NewProviderFactory(),
		config:    cfg,
		stateKeys: make(map[string]StateInfo),
	}

	// Register supported providers
	service.registerProviders()

	return service
}

// registerProviders registers providers
func (s *OAuthService) registerProviders() {
	// Google provider registration
	if s.config.Google.ClientID != "" {
		googleProvider := NewGoogleProvider(s.config.Google)
		s.factory.Register("google", googleProvider)
	}
}

// GetAuthURL generates authentication URL for login
func (s *OAuthService) GetAuthURLForLogin(providerName string) (string, error) {
	provider, err := s.factory.Get(providerName)
	if err != nil {
		return "", err
	}

	if err := provider.ValidateConfig(); err != nil {
		return "", fmt.Errorf("provider config invalid: %w", err)
	}

	// Generate state token for login
	state := s.generateState(nil, providerName, "login")

	// Generate authentication URL
	authURL := provider.GetAuthURL(state)

	return authURL, nil
}

// GetAuthURL generates authentication URL for connecting account (requires user ID)
func (s *OAuthService) GetAuthURL(providerName string, userID uuid.UUID, action string) (string, error) {
	provider, err := s.factory.Get(providerName)
	if err != nil {
		return "", err
	}

	if err := provider.ValidateConfig(); err != nil {
		return "", fmt.Errorf("provider config invalid: %w", err)
	}

	// Generate state token
	state := s.generateState(&userID, providerName, action)

	// Generate authentication URL
	authURL := provider.GetAuthURL(state)

	return authURL, nil
}

// HandleCallback handles OAuth callback
func (s *OAuthService) HandleCallback(ctx context.Context, providerName, code, state string) (*CallbackResult, error) {
	// Validate state token
	stateInfo, err := s.validateState(state)
	if err != nil {
		return nil, fmt.Errorf("invalid state: %w", err)
	}

	// Get provider
	provider, err := s.factory.Get(providerName)
	if err != nil {
		return nil, err
	}

	// Exchange authorization code for access token
	tokenResp, err := provider.ExchangeCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	// Get user profile information
	profile, err := provider.GetUserProfile(ctx, tokenResp.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}

	result := &CallbackResult{
		UserID:    stateInfo.UserID,
		Provider:  providerName,
		Action:    stateInfo.Action,
		Profile:   profile,
		TokenInfo: tokenResp,
		IsNewUser: false,
	}

	// Delete state token (completed)
	delete(s.stateKeys, state)

	return result, nil
}

// CallbackResult callback processing result
type CallbackResult struct {
	UserID    *uuid.UUID    `json:"user_id,omitempty"`
	Provider  string        `json:"provider"`
	Action    string        `json:"action"`
	Profile   *UserProfile  `json:"profile"`
	TokenInfo *TokenResponse `json:"token_info"`
	IsNewUser bool          `json:"is_new_user"`
}

// GetSupportedProviders returns list of supported providers
func (s *OAuthService) GetSupportedProviders() []string {
	return s.factory.GetSupportedProviders()
}

// generateState generates state token
func (s *OAuthService) generateState(userID *uuid.UUID, provider, action string) string {
	// Generate random string
	randomBytes := make([]byte, 16)
	rand.Read(randomBytes)
	randomStr := base64.URLEncoding.EncodeToString(randomBytes)

	// Generate state information
	now := time.Now()
	stateInfo := StateInfo{
		UserID:    userID,
		Provider:  provider,
		Action:    action,
		CreatedAt: now,
		ExpiresAt: now.Add(10 * time.Minute), // Expires in 10 minutes
	}

	// State token format: userID:provider:action:timestamp:random
	var userIDStr string
	if userID != nil {
		userIDStr = userID.String()
	} else {
		userIDStr = "nil"
	}
	
	stateToken := fmt.Sprintf("%s:%s:%s:%d:%s",
		userIDStr, provider, action, now.Unix(), randomStr)

	// Base64 encoding
	encodedState := base64.URLEncoding.EncodeToString([]byte(stateToken))

	// Store in memory (Redis recommended in production)
	s.stateKeys[encodedState] = stateInfo

	return encodedState
}

// validateState validates state token
func (s *OAuthService) validateState(stateToken string) (*StateInfo, error) {
	// Lookup state token
	stateInfo, exists := s.stateKeys[stateToken]
	if !exists {
		return nil, fmt.Errorf("state token not found")
	}

	// Check expiration time
	if time.Now().After(stateInfo.ExpiresAt) {
		delete(s.stateKeys, stateToken) // Delete expired token
		return nil, fmt.Errorf("state token expired")
	}

	// Decode and validate token
	decodedBytes, err := base64.URLEncoding.DecodeString(stateToken)
	if err != nil {
		return nil, fmt.Errorf("invalid state format")
	}

	parts := strings.Split(string(decodedBytes), ":")
	if len(parts) != 5 {
		return nil, fmt.Errorf("invalid state format")
	}

	// Validate userID
	if parts[0] != "nil" {
		userID, err := uuid.Parse(parts[0])
		if err != nil {
			return nil, fmt.Errorf("invalid user ID in state")
		}

		if stateInfo.UserID == nil || userID != *stateInfo.UserID {
			return nil, fmt.Errorf("user ID mismatch")
		}
	}

	// Validate provider and action
	if parts[1] != stateInfo.Provider || parts[2] != stateInfo.Action {
		return nil, fmt.Errorf("provider or action mismatch")
	}

	return &stateInfo, nil
}

// CleanupExpiredStates cleans up expired state tokens (call periodically)
func (s *OAuthService) CleanupExpiredStates() {
	now := time.Now()
	for token, stateInfo := range s.stateKeys {
		if now.After(stateInfo.ExpiresAt) {
			delete(s.stateKeys, token)
		}
	}
}
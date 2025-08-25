package oauth

import (
	"context"
	"fmt"
	"net/url"
)

// OAuthProvider interface - all OAuth providers must implement these methods
type OAuthProvider interface {
	// GetAuthURL generates authentication URL
	GetAuthURL(state string) string

	// ExchangeCode exchanges authorization code for access token
	ExchangeCode(ctx context.Context, code string) (*TokenResponse, error)

	// GetUserProfile gets user profile information using access token
	GetUserProfile(ctx context.Context, accessToken string) (*UserProfile, error)

	// GetProviderName returns provider name
	GetProviderName() string

	// ValidateConfig validates configuration
	ValidateConfig() error
}

// TokenResponse OAuth token exchange response
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Scope        string `json:"scope,omitempty"`
}

// UserProfile user profile information (standardized format)
type UserProfile struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	DisplayName string `json:"display_name"`
	ProfileURL  string `json:"profile_url,omitempty"`
	Avatar      string `json:"avatar,omitempty"`
	Provider    string `json:"provider"`
	RawData     map[string]interface{} `json:"raw_data,omitempty"`
}

// ProviderFactory OAuth provider factory
type ProviderFactory struct {
	providers map[string]OAuthProvider
}

// NewProviderFactory creates factory
func NewProviderFactory() *ProviderFactory {
	return &ProviderFactory{
		providers: make(map[string]OAuthProvider),
	}
}

// Register registers provider
func (f *ProviderFactory) Register(name string, provider OAuthProvider) {
	f.providers[name] = provider
}

// Get retrieves provider
func (f *ProviderFactory) Get(name string) (OAuthProvider, error) {
	provider, exists := f.providers[name]
	if !exists {
		return nil, fmt.Errorf("oauth provider '%s' not found", name)
	}
	return provider, nil
}

// GetSupportedProviders returns list of supported providers
func (f *ProviderFactory) GetSupportedProviders() []string {
	var names []string
	for name := range f.providers {
		names = append(names, name)
	}
	return names
}

// Helper functions

// BuildURL URL helper
func BuildURL(baseURL string, params map[string]string) string {
	u, _ := url.Parse(baseURL)
	q := u.Query()
	for key, value := range params {
		q.Set(key, value)
	}
	u.RawQuery = q.Encode()
	return u.String()
}
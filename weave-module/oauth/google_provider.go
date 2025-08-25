package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"weave-module/config"
)

// GoogleProvider Google OAuth provider
type GoogleProvider struct {
	config config.GoogleOAuthConfig
	client *http.Client
}

// NewGoogleProvider creates Google provider
func NewGoogleProvider(cfg config.GoogleOAuthConfig) *GoogleProvider {
	return &GoogleProvider{
		config: cfg,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetProviderName returns provider name
func (p *GoogleProvider) GetProviderName() string {
	return "google"
}

// ValidateConfig validates configuration
func (p *GoogleProvider) ValidateConfig() error {
	if p.config.ClientID == "" {
		return fmt.Errorf("google client_id is required")
	}
	if p.config.ClientSecret == "" {
		return fmt.Errorf("google client_secret is required")
	}
	if p.config.RedirectURL == "" {
		return fmt.Errorf("google redirect_url is required")
	}
	return nil
}

// GetAuthURL generates authentication URL
func (p *GoogleProvider) GetAuthURL(state string) string {
	scopes := p.config.Scopes
	if scopes == "" {
		scopes = "profile email"
	}

	params := map[string]string{
		"response_type": "code",
		"client_id":     p.config.ClientID,
		"redirect_uri":  p.config.RedirectURL,
		"scope":         scopes,
		"state":         state,
		"access_type":   "offline",
		"prompt":        "consent",
	}

	return BuildURL("https://accounts.google.com/o/oauth2/v2/auth", params)
}

// ExchangeCode exchanges authorization code for access token
func (p *GoogleProvider) ExchangeCode(ctx context.Context, code string) (*TokenResponse, error) {
	tokenURL := "https://oauth2.googleapis.com/token"

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", p.config.RedirectURL)
	data.Set("client_id", p.config.ClientID)
	data.Set("client_secret", p.config.ClientSecret)

	req, err := http.NewRequestWithContext(ctx, "POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("google token exchange failed %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("failed to parse token response: %w", err)
	}

	return &tokenResp, nil
}

// GetUserProfile gets user profile information using access token
func (p *GoogleProvider) GetUserProfile(ctx context.Context, accessToken string) (*UserProfile, error) {
	profileURL := "https://www.googleapis.com/oauth2/v2/userinfo"

	req, err := http.NewRequestWithContext(ctx, "GET", profileURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create profile request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch profile: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("google profile API error %d: %s", resp.StatusCode, string(body))
	}

	var googleProfile GoogleProfile
	if err := json.NewDecoder(resp.Body).Decode(&googleProfile); err != nil {
		return nil, fmt.Errorf("failed to parse profile response: %w", err)
	}

	userProfile := &UserProfile{
		ID:          googleProfile.ID,
		Email:       googleProfile.Email,
		FirstName:   googleProfile.GivenName,
		LastName:    googleProfile.FamilyName,
		DisplayName: googleProfile.Name,
		Avatar:      googleProfile.Picture,
		Provider:    "google",
		RawData:     map[string]interface{}{"google_profile": googleProfile},
	}

	return userProfile, nil
}

// GoogleProfile Google API response structure
type GoogleProfile struct {
	ID         string `json:"id"`
	Email      string `json:"email"`
	Name       string `json:"name"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
	Picture    string `json:"picture"`
	Locale     string `json:"locale"`
}
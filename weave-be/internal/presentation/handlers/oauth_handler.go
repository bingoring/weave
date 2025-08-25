package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"weave-module/config"
	"weave-module/oauth"
	"weave-module/utils"
	"weave-be/internal/application/services"
)

// OAuthHandler OAuth related handler
type OAuthHandler struct {
	userService  *services.UserApplicationService
	oauthService *oauth.OAuthService
	config       *config.Config
}

// NewOAuthHandler creates OAuth handler
func NewOAuthHandler(userService *services.UserApplicationService, cfg *config.Config) *OAuthHandler {
	return &OAuthHandler{
		userService:  userService,
		oauthService: oauth.NewOAuthService(cfg.OAuth),
		config:       cfg,
	}
}

// StartGoogleLogin starts Google OAuth login flow
// GET /api/auth/google/login
func (h *OAuthHandler) StartGoogleLogin(c *gin.Context) {
	// Generate OAuth authentication URL for login
	authURL, err := h.oauthService.GetAuthURLForLogin("google")
	if err != nil {
		if err.Error() == "oauth provider 'google' not found" {
			utils.ErrorResponse(c, fmt.Errorf("Google OAuth is not configured"))
			return
		}
		utils.ErrorResponse(c, fmt.Errorf("Failed to generate Google OAuth URL"))
		return
	}

	utils.SuccessResponse(c, "Google OAuth URL generated successfully", gin.H{
		"auth_url": authURL,
		"provider": "google",
		"action":   "login",
	})
}

// StartGoogleConnect starts Google OAuth connect flow (requires authentication)
// GET /api/auth/google/connect
func (h *OAuthHandler) StartGoogleConnect(c *gin.Context) {
	// User authentication check (connect feature is only available for logged-in users)
	userIDValue, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, fmt.Errorf("User authentication required for Google account connection"))
		return
	}

	userID, ok := userIDValue.(uuid.UUID)
	if !ok {
		utils.ErrorResponse(c, fmt.Errorf("Invalid user ID format"))
		return
	}

	// Generate OAuth authentication URL for connection
	authURL, err := h.oauthService.GetAuthURL("google", userID, "connect")
	if err != nil {
		if err.Error() == "oauth provider 'google' not found" {
			utils.ErrorResponse(c, fmt.Errorf("Google OAuth is not configured"))
			return
		}
		utils.ErrorResponse(c, fmt.Errorf("Failed to generate Google OAuth URL"))
		return
	}

	utils.SuccessResponse(c, "Google OAuth URL generated successfully", gin.H{
		"auth_url": authURL,
		"provider": "google",
		"action":   "connect",
	})
}

// GoogleOAuthCallback handles Google OAuth callback
// GET /api/auth/google/callback
func (h *OAuthHandler) GoogleOAuthCallback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")
	errorParam := c.Query("error")

	// Error handling
	if errorParam != "" {
		errorDescription := c.Query("error_description")
		redirectURL := fmt.Sprintf("%s/login?error=%s&description=%s",
			h.config.App.Name, errorParam, errorDescription)
		c.Redirect(http.StatusFound, redirectURL)
		return
	}

	// Code validation
	if code == "" {
		redirectURL := fmt.Sprintf("%s/login?error=no_code", h.config.App.Name)
		c.Redirect(http.StatusFound, redirectURL)
		return
	}

	// Handle OAuth callback with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Determine if this is login or connect based on state
	result, err := h.oauthService.HandleCallback(ctx, "google", code, state)
	if err != nil {
		redirectURL := fmt.Sprintf("%s/login?error=oauth_failed&provider=google",
			h.config.App.Name)
		c.Redirect(http.StatusFound, redirectURL)
		return
	}

	// Handle login action
	if result.Action == "login" {
		loginResponse, err := h.userService.GoogleOAuthLogin(ctx, code, state)
		if err != nil {
			redirectURL := fmt.Sprintf("%s/login?error=login_failed&provider=google",
				h.config.App.Name)
			c.Redirect(http.StatusFound, redirectURL)
			return
		}

		// Success redirect for login
		redirectURL := fmt.Sprintf("%s/dashboard?token=%s&user=%s",
			h.config.App.Name, loginResponse.Token, loginResponse.User.Username)
		c.Redirect(http.StatusFound, redirectURL)
		return
	}

	// Handle connect action
	if result.Action == "connect" {
		if result.UserID == nil {
			redirectURL := fmt.Sprintf("%s/settings?error=invalid_state", h.config.App.Name)
			c.Redirect(http.StatusFound, redirectURL)
			return
		}

		_, err = h.userService.GoogleOAuthConnect(ctx, *result.UserID, code, state)
		if err != nil {
			redirectURL := fmt.Sprintf("%s/settings?error=connection_failed&provider=google",
				h.config.App.Name)
			c.Redirect(http.StatusFound, redirectURL)
			return
		}

		// Success redirect for connect
		redirectURL := fmt.Sprintf("%s/settings?connected=google&name=%s",
			h.config.App.Name, result.Profile.DisplayName)
		c.Redirect(http.StatusFound, redirectURL)
		return
	}

	// Unsupported action
	redirectURL := fmt.Sprintf("%s/login?error=unsupported_action", h.config.App.Name)
	c.Redirect(http.StatusFound, redirectURL)
}

// GetSupportedProviders returns list of supported OAuth providers
// GET /api/auth/providers
func (h *OAuthHandler) GetSupportedProviders(c *gin.Context) {
	providers := h.oauthService.GetSupportedProviders()

	utils.SuccessResponse(c, "Supported OAuth providers retrieved successfully", gin.H{
		"providers": providers,
		"count":     len(providers),
	})
}
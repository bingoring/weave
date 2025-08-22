package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"weave-module/utils"
	"weave-module/errors"
	"weave-be/internal/application/dto"
	"weave-be/internal/application/services"
)

// UserHandler handles HTTP requests related to users
// This follows the Controller pattern and handles presentation concerns
type UserHandler struct {
	userAppService services.UserApplicationService
}

// NewUserHandler creates a new user handler
func NewUserHandler(userAppService services.UserApplicationService) *UserHandler {
	return &UserHandler{
		userAppService: userAppService,
	}
}

// RegisterUser handles user registration
// POST /api/auth/register
func (h *UserHandler) RegisterUser(c *gin.Context) {
	var req dto.RegisterUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	user, err := h.userAppService.RegisterUser(c.Request.Context(), req)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.CreatedResponse(c, "User registered successfully", user)
}

// LoginUser handles user authentication
// POST /api/auth/login
func (h *UserHandler) LoginUser(c *gin.Context) {
	var req dto.LoginUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	response, err := h.userAppService.LoginUser(c.Request.Context(), req)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, "Login successful", response)
}

// GetProfile gets the authenticated user's profile
// GET /api/users/profile
func (h *UserHandler) GetProfile(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, errors.Unauthorized("User not authenticated"))
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		utils.ErrorResponse(c, errors.BadRequest("Invalid user ID"))
		return
	}

	profile, err := h.userAppService.GetUserProfile(c.Request.Context(), userID)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, "Profile retrieved successfully", profile)
}

// GetUserByID gets a user by ID
// GET /api/users/:id
func (h *UserHandler) GetUserByID(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		utils.ErrorResponse(c, errors.BadRequest("Invalid user ID"))
		return
	}

	user, err := h.userAppService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, "User retrieved successfully", user)
}

// UpdateProfile updates the authenticated user's profile
// PUT /api/users/profile
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, errors.Unauthorized("User not authenticated"))
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		utils.ErrorResponse(c, errors.BadRequest("Invalid user ID"))
		return
	}

	var req dto.UpdateUserProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	user, err := h.userAppService.UpdateUserProfile(c.Request.Context(), userID, req)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, "Profile updated successfully", user)
}

// FollowUser follows another user
// POST /api/users/:id/follow
func (h *UserHandler) FollowUser(c *gin.Context) {
	followerIDStr, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, errors.Unauthorized("User not authenticated"))
		return
	}

	followerID, err := uuid.Parse(followerIDStr.(string))
	if err != nil {
		utils.ErrorResponse(c, errors.BadRequest("Invalid follower ID"))
		return
	}

	followingIDStr := c.Param("id")
	followingID, err := uuid.Parse(followingIDStr)
	if err != nil {
		utils.ErrorResponse(c, errors.BadRequest("Invalid user ID to follow"))
		return
	}

	err = h.userAppService.FollowUser(c.Request.Context(), followerID, followingID)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, "User followed successfully", nil)
}

// UnfollowUser unfollows another user
// DELETE /api/users/:id/follow
func (h *UserHandler) UnfollowUser(c *gin.Context) {
	followerIDStr, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, errors.Unauthorized("User not authenticated"))
		return
	}

	followerID, err := uuid.Parse(followerIDStr.(string))
	if err != nil {
		utils.ErrorResponse(c, errors.BadRequest("Invalid follower ID"))
		return
	}

	followingIDStr := c.Param("id")
	followingID, err := uuid.Parse(followingIDStr)
	if err != nil {
		utils.ErrorResponse(c, errors.BadRequest("Invalid user ID to unfollow"))
		return
	}

	err = h.userAppService.UnfollowUser(c.Request.Context(), followerID, followingID)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, "User unfollowed successfully", nil)
}

// GetFollowers gets user's followers
// GET /api/users/:id/followers
func (h *UserHandler) GetFollowers(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		utils.ErrorResponse(c, errors.BadRequest("Invalid user ID"))
		return
	}

	page, limit := h.getPaginationParams(c)

	followers, err := h.userAppService.GetFollowers(c.Request.Context(), userID, page, limit)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	pagination := utils.CalculatePagination(page, limit, int64(followers.Total))
	utils.PaginatedSuccessResponse(c, "Followers retrieved successfully", followers.Users, pagination)
}

// GetFollowing gets users that a user is following
// GET /api/users/:id/following
func (h *UserHandler) GetFollowing(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		utils.ErrorResponse(c, errors.BadRequest("Invalid user ID"))
		return
	}

	page, limit := h.getPaginationParams(c)

	following, err := h.userAppService.GetFollowing(c.Request.Context(), userID, page, limit)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	pagination := utils.CalculatePagination(page, limit, int64(following.Total))
	utils.PaginatedSuccessResponse(c, "Following retrieved successfully", following.Users, pagination)
}

// SearchUsers searches for users
// GET /api/users/search
func (h *UserHandler) SearchUsers(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		utils.ErrorResponse(c, errors.BadRequest("Search query is required"))
		return
	}

	page, limit := h.getPaginationParams(c)

	results, err := h.userAppService.SearchUsers(c.Request.Context(), query, page, limit)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	pagination := utils.CalculatePagination(page, limit, int64(results.Total))
	utils.PaginatedSuccessResponse(c, "Users found", results.Users, pagination)
}

// Helper methods
func (h *UserHandler) getPaginationParams(c *gin.Context) (page, limit int) {
	page = 1
	limit = 20

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	return page, limit
}
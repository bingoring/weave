package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"weave-module/utils"
	"weave-module/errors"
	"weave-be/internal/application/dto"
	"weave-be/internal/application/services"
)

// UserHandler handles HTTP requests related to users using Use Case based architecture
type UserHandler struct {
	userService *services.UserApplicationService
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService *services.UserApplicationService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// SendEmailVerification handles email verification code sending
func (h *UserHandler) SendEmailVerification(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, errors.BadRequest("Invalid request body"))
		return
	}

	response, err := h.userService.SendEmailVerification(c.Request.Context(), req.Email)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, "Verification code sent", response)
}

// VerifyEmailAuth handles email verification code confirmation and login
func (h *UserHandler) VerifyEmailAuth(c *gin.Context) {
	var req struct {
		Code string `json:"code" binding:"required,len=6"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, errors.BadRequest("Invalid request body"))
		return
	}

	response, err := h.userService.VerifyEmailAuth(c.Request.Context(), req.Code)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, "Email verification successful", response)
}

// GetProfile handles get user profile requests (requires authentication)
func (h *UserHandler) GetProfile(c *gin.Context) {
	userIDValue, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, errors.Unauthorized("User not authenticated"))
		return
	}

	userID, ok := userIDValue.(uuid.UUID)
	if !ok {
		utils.ErrorResponse(c, errors.InternalServerError("Invalid user ID format"))
		return
	}

	profile, err := h.userService.GetUserProfile(c.Request.Context(), userID)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, "Profile retrieved successfully", profile)
}

// GetUserByID handles get user by ID requests
func (h *UserHandler) GetUserByID(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		utils.ErrorResponse(c, errors.BadRequest("Invalid user ID"))
		return
	}

	user, err := h.userService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, "User retrieved successfully", user)
}

// UpdateProfile handles update user profile requests
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userIDValue, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, errors.Unauthorized("User not authenticated"))
		return
	}

	userID, ok := userIDValue.(uuid.UUID)
	if !ok {
		utils.ErrorResponse(c, errors.InternalServerError("Invalid user ID format"))
		return
	}

	var req dto.UpdateUserProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, errors.BadRequest("Invalid request body"))
		return
	}

	user, err := h.userService.UpdateUserProfile(c.Request.Context(), userID, req)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, "Profile updated successfully", user)
}

// FollowUser handles follow user requests
func (h *UserHandler) FollowUser(c *gin.Context) {
	followerIDValue, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, errors.Unauthorized("User not authenticated"))
		return
	}

	followerID, ok := followerIDValue.(uuid.UUID)
	if !ok {
		utils.ErrorResponse(c, errors.InternalServerError("Invalid user ID format"))
		return
	}

	followingIDStr := c.Param("id")
	followingID, err := uuid.Parse(followingIDStr)
	if err != nil {
		utils.ErrorResponse(c, errors.BadRequest("Invalid user ID"))
		return
	}

	err = h.userService.FollowUser(c.Request.Context(), followerID, followingID)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, "User followed successfully", nil)
}

// UnfollowUser handles unfollow user requests
func (h *UserHandler) UnfollowUser(c *gin.Context) {
	followerIDValue, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, errors.Unauthorized("User not authenticated"))
		return
	}

	followerID, ok := followerIDValue.(uuid.UUID)
	if !ok {
		utils.ErrorResponse(c, errors.InternalServerError("Invalid user ID format"))
		return
	}

	followingIDStr := c.Param("id")
	followingID, err := uuid.Parse(followingIDStr)
	if err != nil {
		utils.ErrorResponse(c, errors.BadRequest("Invalid user ID"))
		return
	}

	err = h.userService.UnfollowUser(c.Request.Context(), followerID, followingID)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, "User unfollowed successfully", nil)
}

// GetFollowers handles get followers requests
func (h *UserHandler) GetFollowers(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		utils.ErrorResponse(c, errors.BadRequest("Invalid user ID"))
		return
	}

	page, limit := utils.GetPaginationParams(c)

	followers, err := h.userService.GetFollowers(c.Request.Context(), userID, page, limit)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	pagination := utils.CalculatePagination(page, limit, int64(followers.Total))
	utils.PaginatedSuccessResponse(c, "Followers retrieved successfully", followers.Users, pagination)
}

// GetFollowing handles get following requests
func (h *UserHandler) GetFollowing(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		utils.ErrorResponse(c, errors.BadRequest("Invalid user ID"))
		return
	}

	page, limit := utils.GetPaginationParams(c)

	following, err := h.userService.GetFollowing(c.Request.Context(), userID, page, limit)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	pagination := utils.CalculatePagination(page, limit, int64(following.Total))
	utils.PaginatedSuccessResponse(c, "Following retrieved successfully", following.Users, pagination)
}

// SearchUsers handles search users requests
func (h *UserHandler) SearchUsers(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		utils.ErrorResponse(c, errors.BadRequest("Search query is required"))
		return
	}

	page, limit := utils.GetPaginationParams(c)

	users, err := h.userService.SearchUsers(c.Request.Context(), query, page, limit)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	pagination := utils.CalculatePagination(page, limit, int64(users.Total))
	utils.PaginatedSuccessResponse(c, "Users retrieved successfully", users.Users, pagination)
}
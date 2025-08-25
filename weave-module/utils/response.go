package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"weave-module/errors"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

type PaginatedResponse struct {
	Success    bool        `json:"success"`
	Message    string      `json:"message,omitempty"`
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
}


func SuccessResponse(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func CreatedResponse(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func PaginatedSuccessResponse(c *gin.Context, message string, data interface{}, pagination Pagination) {
	c.JSON(http.StatusOK, PaginatedResponse{
		Success:    true,
		Message:    message,
		Data:       data,
		Pagination: pagination,
	})
}

func ErrorResponse(c *gin.Context, err error) {
	if appErr, ok := err.(*errors.AppError); ok {
		c.JSON(appErr.Code, Response{
			Success: false,
			Error:   appErr.Message,
		})
		return
	}

	// Default to internal server error for unknown errors
	c.JSON(http.StatusInternalServerError, Response{
		Success: false,
		Error:   "Internal server error",
	})
}

func ValidationErrorResponse(c *gin.Context, validationErrors interface{}) {
	c.JSON(http.StatusBadRequest, Response{
		Success: false,
		Error:   "Validation failed",
		Data:    validationErrors,
	})
}
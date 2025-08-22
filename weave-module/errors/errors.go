package errors

import (
	"fmt"
	"net/http"
)

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(code int, message, details string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// Common error constructors
func BadRequest(message string) *AppError {
	return NewAppError(http.StatusBadRequest, message, "")
}

func BadRequestWithDetails(message, details string) *AppError {
	return NewAppError(http.StatusBadRequest, message, details)
}

func Unauthorized(message string) *AppError {
	return NewAppError(http.StatusUnauthorized, message, "")
}

func Forbidden(message string) *AppError {
	return NewAppError(http.StatusForbidden, message, "")
}

func NotFound(message string) *AppError {
	return NewAppError(http.StatusNotFound, message, "")
}

func Conflict(message string) *AppError {
	return NewAppError(http.StatusConflict, message, "")
}

func InternalServerError(message string) *AppError {
	return NewAppError(http.StatusInternalServerError, message, "")
}

func InternalServerErrorWithDetails(message, details string) *AppError {
	return NewAppError(http.StatusInternalServerError, message, details)
}

func ValidationError(field, message string) *AppError {
	return NewAppError(http.StatusBadRequest, fmt.Sprintf("Validation failed for field '%s': %s", field, message), "")
}

// Specific business logic errors
var (
	ErrUserNotFound         = NotFound("User not found")
	ErrUserAlreadyExists    = Conflict("User already exists")
	ErrInvalidCredentials   = Unauthorized("Invalid credentials")
	ErrInvalidToken         = Unauthorized("Invalid or expired token")
	ErrInsufficientPermissions = Forbidden("Insufficient permissions")
	
	ErrWeaveNotFound        = NotFound("Weave not found")
	ErrWeaveAlreadyExists   = Conflict("Weave already exists")
	ErrWeaveNotPublished    = BadRequest("Weave is not published")
	
	ErrChannelNotFound      = NotFound("Channel not found")
	ErrChannelAlreadyExists = Conflict("Channel already exists")
	
	ErrInvalidInput         = BadRequest("Invalid input provided")
	ErrMissingRequiredField = BadRequest("Missing required field")
	
	ErrDatabaseConnection   = InternalServerError("Database connection error")
	ErrCacheConnection      = InternalServerError("Cache connection error")
	ErrQueueConnection      = InternalServerError("Queue connection error")
)

// Helper functions to check error types
func IsNotFound(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == http.StatusNotFound
	}
	return false
}

func IsBadRequest(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == http.StatusBadRequest
	}
	return false
}

func IsUnauthorized(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == http.StatusUnauthorized
	}
	return false
}

func IsConflict(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == http.StatusConflict
	}
	return false
}

func IsInternalServerError(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == http.StatusInternalServerError
	}
	return false
}
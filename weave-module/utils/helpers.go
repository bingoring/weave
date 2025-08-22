package utils

import (
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// String helpers
func TrimAndLower(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func ToSlug(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "_", "-")
	// Remove any character that's not alphanumeric or hyphen
	var result strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result.WriteRune(r)
		}
	}
	return result.String()
}

func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// UUID helpers
func ParseUUID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}

func IsValidUUID(s string) bool {
	_, err := uuid.Parse(s)
	return err == nil
}

// Pagination helpers
func GetPaginationParams(c *gin.Context) (page, limit int) {
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

func CalculatePagination(page, limit int, total int64) Pagination {
	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	hasNext := page < totalPages
	hasPrev := page > 1

	return Pagination{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
		HasNext:    hasNext,
		HasPrev:    hasPrev,
	}
}

func GetOffset(page, limit int) int {
	return (page - 1) * limit
}

// Time helpers
func FormatTimeAgo(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	if diff < time.Minute {
		return "just now"
	} else if diff < time.Hour {
		minutes := int(diff.Minutes())
		if minutes == 1 {
			return "1 minute ago"
		}
		return strconv.Itoa(minutes) + " minutes ago"
	} else if diff < 24*time.Hour {
		hours := int(diff.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return strconv.Itoa(hours) + " hours ago"
	} else if diff < 30*24*time.Hour {
		days := int(diff.Hours() / 24)
		if days == 1 {
			return "1 day ago"
		}
		return strconv.Itoa(days) + " days ago"
	} else if diff < 365*24*time.Hour {
		months := int(diff.Hours() / (24 * 30))
		if months == 1 {
			return "1 month ago"
		}
		return strconv.Itoa(months) + " months ago"
	} else {
		years := int(diff.Hours() / (24 * 365))
		if years == 1 {
			return "1 year ago"
		}
		return strconv.Itoa(years) + " years ago"
	}
}

// File upload helpers
func GetAllowedImageExtensions() []string {
	return []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
}

func GetAllowedVideoExtensions() []string {
	return []string{".mp4", ".mov", ".avi", ".mkv", ".webm"}
}

func IsImageFile(filename string) bool {
	ext := strings.ToLower(filename[strings.LastIndex(filename, "."):])
	return Contains(GetAllowedImageExtensions(), ext)
}

func IsVideoFile(filename string) bool {
	ext := strings.ToLower(filename[strings.LastIndex(filename, "."):])
	return Contains(GetAllowedVideoExtensions(), ext)
}

// Search helpers
func GetSearchQuery(c *gin.Context) string {
	return strings.TrimSpace(c.Query("q"))
}

func GetSortParams(c *gin.Context) (string, string) {
	sortBy := c.Query("sort_by")
	if sortBy == "" {
		sortBy = "created_at"
	}

	order := c.Query("order")
	if order != "asc" && order != "desc" {
		order = "desc"
	}

	return sortBy, order
}

// Filter helpers
func GetDateRangeParams(c *gin.Context) (time.Time, time.Time, error) {
	var startDate, endDate time.Time
	var err error

	if start := c.Query("start_date"); start != "" {
		startDate, err = time.Parse("2006-01-02", start)
		if err != nil {
			return startDate, endDate, err
		}
	}

	if end := c.Query("end_date"); end != "" {
		endDate, err = time.Parse("2006-01-02", end)
		if err != nil {
			return startDate, endDate, err
		}
		// Set to end of day
		endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
	}

	return startDate, endDate, nil
}
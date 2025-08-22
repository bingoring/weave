package models

import (
	"github.com/google/uuid"
)

// ParseUUID is a helper function to parse UUID from string
func ParseUUID(s string) uuid.UUID {
	id, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil
	}
	return id
}
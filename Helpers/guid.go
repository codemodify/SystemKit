package Helpers

import (
	"github.com/google/uuid"
)

// NewGUID -
func NewGUID() string {
	guidAsString := RandomString(50)
	id, err := uuid.NewUUID()
	if err == nil {
		guidAsString = id.String()
	}

	return guidAsString
}

// NewGUIDWithLength -
func NewGUIDWithLength(length int) string {
	guidAsString := RandomString(length)
	id, err := uuid.NewUUID()
	if err == nil {
		guidAsString = id.String()
	}

	return guidAsString
}

package valueobject

import (
	"errors"

	"github.com/google/uuid"
)

type UserID struct {
	value string
}

func NewUserID() UserID {
	// WHY: Universally unique, no coordination needed
	return UserID{value: uuid.New().String()}
}

func NewUserIDFromString(id string) (UserID, error) {
	if id == "" {
		return UserID{}, errors.New("user ID cannot be empty")
	}

	// Validate it's a valid UUID format
	if _, err := uuid.Parse(id); err != nil {
		return UserID{}, errors.New("invalid user ID format")
	}

	return UserID{value: id}, nil
}

// String returns the ID as a string
func (id UserID) String() string {
	return id.value
}

func (id UserID) Equals(other UserID) bool {
	return id.value == other.value
}

func (id UserID) IsEmpty() bool {
	return id.value == ""
}

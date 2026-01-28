package testutil

import (
	"strconv"
	"testing"
	"time"

	"github.com/Ruseigha/LabukaAuth/internal/domain/entity"
	"github.com/Ruseigha/LabukaAuth/internal/domain/valueobject"
)

// CreateTestUser creates a user entity for testing
// WHY: Reusable test data creation
func CreateTestUser(t *testing.T, email, password string) *entity.User {
	t.Helper()

	emailVO, err := valueobject.NewEmail(email)
	if err != nil {
		t.Fatalf("Failed to create test email: %v", err)
	}

	passwordVO, err := valueobject.NewPassword(password)
	if err != nil {
		t.Fatalf("Failed to create test password: %v", err)
	}

	user, err := entity.NewUser(emailVO, passwordVO)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	return user
}

// CreateTestUserWithID creates a user with specific ID (for reconstruction)
func CreateTestUserWithID(t *testing.T, id, email, password string, isActive bool) *entity.User {
	t.Helper()

	userID, err := valueobject.NewUserIDFromString(id)
	if err != nil {
		t.Fatalf("Failed to create test user ID: %v", err)
	}

	emailVO, err := valueobject.NewEmail(email)
	if err != nil {
		t.Fatalf("Failed to create test email: %v", err)
	}

	passwordVO, err := valueobject.NewPassword(password)
	if err != nil {
		t.Fatalf("Failed to create test password: %v", err)
	}

	now := time.Now().UTC()

	return entity.ReconstructUser(
		userID,
		emailVO,
		passwordVO,
		now,
		now,
		isActive,
	)
}

// ValidTestEmail returns a valid test email
func ValidTestEmail() string {
	return "test@example.com"
}

// ValidTestPassword returns a valid test password
func ValidTestPassword() string {
	return "SecureP@ss123"
}

// UniqueEmail generates a unique email for testing
// WHY: Avoid conflicts when inserting multiple users
func UniqueEmail(prefix string) string {
	timestamp := time.Now().UnixNano()
	return prefix + "_" + strconv.FormatInt(timestamp, 10) + "@example.com"
}

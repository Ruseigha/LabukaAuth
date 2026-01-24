package valueobject_test

import (
	"testing"

	"github.com/Ruseigha/LabukaAuth/internal/domain/valueobject"
)

func TestNewUserID(t *testing.T) {
	// Generate multiple IDs to ensure uniqueness
	id1 := valueobject.NewUserID()
	id2 := valueobject.NewUserID()

	// Should not be empty
	if id1.String() == "" {
		t.Error("NewUserID() returned empty ID")
	}

	// Should be unique (extremely high probability)
	if id1.Equals(id2) {
		t.Error("NewUserID() generated duplicate IDs")
	}

	// Should be valid UUID format (36 characters with hyphens)
	if len(id1.String()) != 36 {
		t.Errorf("NewUserID() generated ID with wrong length: %d, want 36", len(id1.String()))
	}
}

func TestNewUserIDFromString_Valid(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "valid UUID v4",
			input: "550e8400-e29b-41d4-a716-446655440000",
		},
		{
			name:  "another valid UUID",
			input: "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := valueobject.NewUserIDFromString(tt.input)

			if err != nil {
				t.Errorf("NewUserIDFromString() unexpected error = %v", err)
				return
			}

			if id.String() != tt.input {
				t.Errorf("NewUserIDFromString() = %v, want %v", id.String(), tt.input)
			}
		})
	}
}

func TestNewUserIDFromString_Invalid(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantError string
	}{
		{
			name:      "empty string",
			input:     "",
			wantError: "user ID cannot be empty",
		},
		{
			name:      "invalid format",
			input:     "not-a-uuid",
			wantError: "invalid user ID format",
		},
		{
			name:      "random string",
			input:     "12345",
			wantError: "invalid user ID format",
		},
		{
			name:      "malformed UUID",
			input:     "550e8400-e29b-41d4-a716",
			wantError: "invalid user ID format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := valueobject.NewUserIDFromString(tt.input)

			if err == nil {
				t.Error("NewUserIDFromString() expected error, got nil")
				return
			}

			if err.Error() != tt.wantError {
				t.Errorf("NewUserIDFromString() error = %v, want %v", err.Error(), tt.wantError)
			}
		})
	}
}

func TestUserID_Equals(t *testing.T) {
	id1, _ := valueobject.NewUserIDFromString("550e8400-e29b-41d4-a716-446655440000")
	id2, _ := valueobject.NewUserIDFromString("550e8400-e29b-41d4-a716-446655440000")
	id3, _ := valueobject.NewUserIDFromString("6ba7b810-9dad-11d1-80b4-00c04fd430c8")

	if !id1.Equals(id2) {
		t.Error("Expected IDs with same value to be equal")
	}

	if id1.Equals(id3) {
		t.Error("Expected IDs with different values to not be equal")
	}
}

func BenchmarkNewUserID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		valueobject.NewUserID()
	}
}
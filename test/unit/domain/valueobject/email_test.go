package valueobject_test

import (
	"testing"

	"github.com/Ruseigha/LabukaAuth/internal/domain/valueobject"
)

func TestNewEmail_ValidEmail(t *testing.T) {
	tests := []struct {
        name  string  // Test case name
        input string  // Input email
        want  string  // Expected normalized email
    }{
        {
            name:  "simple email",
            input: "user@example.com",
            want:  "user@example.com",
        },
        {
            name:  "email with uppercase (should normalize to lowercase)",
            input: "User@Example.COM",
            want:  "user@example.com",
        },
        {
            name:  "email with spaces (should trim)",
            input: "  user@example.com  ",
            want:  "user@example.com",
        },
        {
            name:  "email with plus sign",
            input: "user+tag@example.com",
            want:  "user+tag@example.com",
        },
        {
            name:  "email with dots",
            input: "first.last@example.com",
            want:  "first.last@example.com",
        },
        {
            name:  "email with subdomain",
            input: "user@mail.example.com",
            want:  "user@mail.example.com",
        },
        {
            name:  "email with numbers",
            input: "user123@example.com",
            want:  "user123@example.com",
        },
    }

		 // Run all test cases
		 for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Act - perform the action
            got, err := valueobject.NewEmail(tt.input)
            
            // Assert - check results
            if err != nil {
                t.Errorf("NewEmail() unexpected error = %v", err)
                return
            }
            
            if got.String() != tt.want {
                t.Errorf("NewEmail() = %v, want %v", got.String(), tt.want)
            }
        })
    }
}
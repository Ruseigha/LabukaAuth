package valueobject_test

import (
	"testing"

	"github.com/Ruseigha/LabukaAuth/internal/domain/valueobject"
)

// Valid Email Test
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


// Invalid Emails test
func TestNewEmail_InvalidEmail(t *testing.T) {
    tests := []struct {
        name      string
        input     string
        wantError string  // Expected error message (partial match)
    }{
        {
            name:      "empty email",
            input:     "",
            wantError: "email cannot be empty",
        },
        {
            name:      "whitespace only",
            input:     "   ",
            wantError: "email cannot be empty",
        },
        {
            name:      "missing @",
            input:     "userexample.com",
            wantError: "invalid email format",
        },
        {
            name:      "missing domain",
            input:     "user@",
            wantError: "invalid email format",
        },
        {
            name:      "missing local part",
            input:     "@example.com",
            wantError: "invalid email format",
        },
        {
            name:      "missing TLD",
            input:     "user@example",
            wantError: "invalid email format",
        },
        {
            name:      "multiple @",
            input:     "user@@example.com",
            wantError: "invalid email format",
        },
        {
            name:      "spaces in email",
            input:     "user name@example.com",
            wantError: "invalid email format",
        },
        {
            name:      "invalid characters",
            input:     "user!#$@example.com",
            wantError: "invalid email format",
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := valueobject.NewEmail(tt.input)
            
            // We EXPECT an error
            if err == nil {
                t.Errorf("NewEmail() expected error, got valid email: %v", got.String())
                return
            }
            
            // Check error message contains expected text
            if err.Error() != tt.wantError {
                t.Errorf("NewEmail() error = %v, want %v", err.Error(), tt.wantError)
            }
        })
    }
}
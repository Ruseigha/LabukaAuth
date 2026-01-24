package valueobject_test

import (
	"testing"

	"github.com/Ruseigha/LabukaAuth/internal/domain/valueobject"
)

func TestNewPassword_ValidPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
	}{
		{
			name:     "password with all character types",
			password: "SecureP@ss123",
		},
		{
			name:     "password with uppercase, lowercase, number",
			password: "Password123",
		},
		{
			name:     "password with uppercase, lowercase, special",
			password: "Pass@word",
		},
		{
			name:     "password with lowercase, number, special",
			password: "pass@123",
		},
		{
			name:     "minimum length password",
			password: "Abcd123!", // Exactly 8 characters
		},
		{
			name:     "long password",
			password: "ThisIsAVeryLongP@ssw0rdWith50Characters12345678",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := valueobject.NewPassword(tt.password)

			if err != nil {
				t.Errorf("NewPassword() unexpected error = %v", err)
				return
			}

			// Password should be stored (even if not hashed yet in this layer)
			if got.Hash() == "" {
				t.Error("NewPassword() returned empty hash")
			}
		})
	}
}



func TestNewPassword_InvalidPassword(t *testing.T) {
    tests := []struct {
        name      string
        password  string
        wantError string
    }{
        {
            name:      "too short",
            password:  "Pass1!",  // Only 6 characters
            wantError: "password must be at least 8 characters",
        },
        {
            name:      "empty password",
            password:  "",
            wantError: "password must be at least 8 characters",
        },
        {
            name:      "too long",
            password:  "ThisPasswordIsWayTooLongAndExceeds72CharactersWhichIsTheMaximumAllowed12345678901234567890",
            wantError: "password must be at most 72 characters",
        },
        {
            name:      "only lowercase",
            password:  "passwordonly",
            wantError: "password must contain at least 3 of: uppercase, lowercase, number, special character",
        },
        {
            name:      "only uppercase",
            password:  "PASSWORDONLY",
            wantError: "password must contain at least 3 of: uppercase, lowercase, number, special character",
        },
        {
            name:      "only numbers",
            password:  "12345678",
            wantError: "password must contain at least 3 of: uppercase, lowercase, number, special character",
        },
        {
            name:      "only lowercase and uppercase",
            password:  "PasswordOnly",
            wantError: "password must contain at least 3 of: uppercase, lowercase, number, special character",
        },
        {
            name:      "only lowercase and numbers",
            password:  "password123",
            wantError: "password must contain at least 3 of: uppercase, lowercase, number, special character",
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            _, err := valueobject.NewPassword(tt.password)
            
            if err == nil {
                t.Error("NewPassword() expected error, got nil")
                return
            }
            
            if err.Error() != tt.wantError {
                t.Errorf("NewPassword() error = %v, want %v", err.Error(), tt.wantError)
            }
        })
    }
}
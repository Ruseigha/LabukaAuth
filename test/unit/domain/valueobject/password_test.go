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
package mocks

import "errors"

// MockPasswordHasher is a mock implementation of PasswordHasher
type MockPasswordHasher struct {
	HashFunc    func(password string) (string, error)
	CompareFunc func(hashedPassword, password string) error

	HashCalls    int
	CompareCalls int
}

// Hash implements security.PasswordHasher
func (m *MockPasswordHasher) Hash(password string) (string, error) {
	m.HashCalls++
	if m.HashFunc != nil {
		return m.HashFunc(password)
	}
	// Default: return hash with prefix (easy to verify in tests)
	return "hashed_" + password, nil
}

// Compare implements security.PasswordHasher
func (m *MockPasswordHasher) Compare(hashedPassword, password string) error {
	m.CompareCalls++
	if m.CompareFunc != nil {
		return m.CompareFunc(hashedPassword, password)
	}
	// Default: simple comparison (for testing only!)
	expectedHash := "hashed_" + password
	if hashedPassword != expectedHash {
		return ErrPasswordMismatch
	}
	return nil
}

// ErrPasswordMismatch is returned when password doesn't match
var ErrPasswordMismatch = errors.New("password mismatch")

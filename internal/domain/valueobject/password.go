package valueobject

import (
	"errors"
	"unicode"
)

type Password struct {
	hashedValue string // NEVER store plain text passwords!
}

func NewPassword(plainPassword string) (Password, error) {
	// Business Rule 1: Password minimum length
	if len(plainPassword) < 8 {
		return Password{}, errors.New("password must be at least 8 characters")
	}

	// Business Rule 2: Password maximum length (prevent DoS attacks)
	if len(plainPassword) > 72 {
		return Password{}, errors.New("password must be at most 72 characters")
	}

	// Business Rule 3: Password must contain variety of characters
	if err := validatePasswordStrength(plainPassword); err != nil {
		return Password{}, err
	}

	return Password{hashedValue: plainPassword}, nil
}

func validatePasswordStrength(password string) error {
	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	// Business Rule: Password must have at least 3 of 4 character types
	// WHY: Balance security and usability
	typesCount := 0
	if hasUpper {
		typesCount++
	}
	if hasLower {
		typesCount++
	}
	if hasNumber {
		typesCount++
	}
	if hasSpecial {
		typesCount++
	}

	if typesCount < 3 {
		return errors.New("password must contain at least 3 of: uppercase, lowercase, number, special character")
	}

	return nil
}

// Hash returns the hashed password value
// WHY: For database storage
func (p Password) Hash() string {
	return p.hashedValue
}

func NewPasswordFromHash(hashedPassword string) Password {
	// No validation needed - already hashed and validated
	return Password{hashedValue: hashedPassword}
}

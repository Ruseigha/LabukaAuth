package errors

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidEmailFormat  = errors.New("invalid email format")
	ErrWeakPassword        = errors.New("password does not meet complexity requirements")
	ErrEmptyUserId         = errors.New("user ID cannot be empty")
	ErrEmptyEmail          = errors.New("email cannot be empty")
	ErrPasswordHashMissing = errors.New("password hash is missing")
	ErrInvalidInput        = errors.New("invalid input")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrForbidden           = errors.New("forbidden")
	ErrConflict            = errors.New("resource conflict")
	ErrNotFound            = errors.New("resource not found")
)

type DomainError struct {
	// Type is the base error type (ErrInvalidInput, etc.)
	Type error

	// Message is a human-readable description
	Message string

	// Field is the specific field that caused the error (optional)
	Field string

	// Err is the underlying error (if any)
	Err error
}

func (e *DomainError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("%s: %s (field: %s)", e.Type.Error(), e.Message, e.Field)
	}
	return fmt.Sprintf("%s: %s", e.Type.Error(), e.Message)
}

// Unwrap returns the underlying error
// WHY: Allows errors.Is() and errors.As() to work
func (e *DomainError) Unwrap() error {
	return e.Err
}

// Is checks if error matches target
// WHY: Allows errors.Is(err, ErrInvalidInput) to work
func (e *DomainError) Is(target error) bool {
	return errors.Is(e.Type, target)
}

// NewInvalidInputError creates an invalid input error
func NewInvalidInputError(message string, field string) *DomainError {
	return &DomainError{
		Type:    ErrInvalidInput,
		Message: message,
		Field:   field,
	}
}

// NewUnauthorizedError creates an unauthorized error
func NewUnauthorizedError(message string) *DomainError {
	return &DomainError{
		Type:    ErrUnauthorized,
		Message: message,
	}
}

// NewForbiddenError creates a forbidden error
func NewForbiddenError(message string) *DomainError {
	return &DomainError{
		Type:    ErrForbidden,
		Message: message,
	}
}

// NewConflictError creates a conflict error
func NewConflictError(message string) *DomainError {
	return &DomainError{
		Type:    ErrConflict,
		Message: message,
	}
}

// NewNotFoundError creates a not found error
func NewNotFoundError(message string) *DomainError {
	return &DomainError{
		Type:    ErrNotFound,
		Message: message,
	}
}

// WrapError wraps an error with additional context
// WHY: Preserve error chain (crucial for debugging)
func WrapError(baseType error, message string, err error) *DomainError {
	return &DomainError{
		Type:    baseType,
		Message: message,
		Err:     err,
	}
}

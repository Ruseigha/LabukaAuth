package repository

import (
	"errors"
	"fmt"
)

var (
	ErrUserNotFound        = errors.New("user not found")
	ErrUserAlreadyExists   = errors.New("user already exists")
	ErrDatabaseConnection  = errors.New("database connection error")
	ErrDatabaseQuery       = errors.New("database query error")
	ErrDatabaseTransaction = errors.New("database transaction error")
	ErrInvalidID           = errors.New("invalid ID")
)

type RepositoryError struct {
	Op   string
	Type error
	Err  error
}

func (e *RepositoryError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %v", e.Op, e.Type.Error(), e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Op, e.Type.Error())
}

// Unwrap returns the underlying error
func (e *RepositoryError) Unwrap() error {
	return e.Err
}

// Is checks if error matches target
func (e *RepositoryError) Is(target error) bool {
	return errors.Is(e.Type, target)
}

// NewUserNotFoundError creates a user not found error
func NewUserNotFoundError(op string) *RepositoryError {
	return &RepositoryError{
		Op:   op,
		Type: ErrUserNotFound,
	}
}

// NewUserAlreadyExistsError creates a duplicate user error
func NewUserAlreadyExistsError(op string, err error) *RepositoryError {
	return &RepositoryError{
		Op:   op,
		Type: ErrUserAlreadyExists,
		Err:  err,
	}
}

// NewDatabaseConnectionError creates a connection error
func NewDatabaseConnectionError(op string, err error) *RepositoryError {
	return &RepositoryError{
		Op:   op,
		Type: ErrDatabaseConnection,
		Err:  err,
	}
}

// NewDatabaseQueryError creates a query error
func NewDatabaseQueryError(op string, err error) *RepositoryError {
	return &RepositoryError{
		Op:   op,
		Type: ErrDatabaseQuery,
		Err:  err,
	}
}

// NewInvalidIDError creates an invalid ID error
func NewInvalidIDError(op string, err error) *RepositoryError {
	return &RepositoryError{
		Op:   op,
		Type: ErrInvalidID,
		Err:  err,
	}
}

// WrapRepositoryError wraps any error with repository context
// WHY: Convert arbitrary errors to RepositoryError
func WrapRepositoryError(op string, baseType error, err error) *RepositoryError {
	return &RepositoryError{
		Op:   op,
		Type: baseType,
		Err:  err,
	}
}

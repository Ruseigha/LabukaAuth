package errors

import "errors"

var (
	ErrInvalidEmailFormat  = errors.New("invalid email format")
	ErrWeakPassword        = errors.New("password does not meet complexity requirements")
	ErrEmptyUserId         = errors.New("user ID cannot be empty")
	ErrEmptyEmail          = errors.New("email cannot be empty")
	ErrPasswordHashMissing = errors.New("password hash is missing")
)

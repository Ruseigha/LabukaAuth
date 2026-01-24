package entity

import (
	"errors"
	"time"

	"github.com/Ruseigha/LabukaAuth/internal/domain/valueobject"
)

type User struct {
	id        valueobject.UserID   // Unique identifier
	email     valueobject.Email    // Email (validated)
	password  valueobject.Password // Password (hashed)
	createdAt time.Time            // When user was created
	updatedAt time.Time            // When user was last updated
	isActive  bool                 // Can user log in?
}

func NewUser(email valueobject.Email, password valueobject.Password) (*User, error) {
	// Validate inputs
	if email.IsEmpty() {
		return nil, errors.New("email is required")
	}

	if password.Hash() == "" {
		return nil, errors.New("password is required")
	}

	now := time.Now().UTC()

	return &User{
		id:        valueobject.NewUserID(), // Generate new ID
		email:     email,
		password:  password,
		createdAt: now,
		updatedAt: now,
		isActive:  true, // Business rule: new users are active by default
	}, nil
}

// ReconstructUser recreates a user from stored data
// WHY: When loading from database, we already have ID, timestamps, etc.
// This bypasses validation since data is already validated
func ReconstructUser(
	id valueobject.UserID,
	email valueobject.Email,
	password valueobject.Password,
	createdAt time.Time,
	updatedAt time.Time,
	isActive bool,
) *User {
	return &User{
		id:        id,
		email:     email,
		password:  password,
		createdAt: createdAt,
		updatedAt: updatedAt,
		isActive:  isActive,
	}
}

// Getters - expose data in a controlled way
// WHY: Encapsulation - can change internal representation without breaking API
func (u *User) ID() valueobject.UserID {
	return u.id
}

func (u *User) Email() valueobject.Email {
	return u.email
}

func (u *User) Password() valueobject.Password {
	return u.password
}

func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

func (u *User) UpdatedAt() time.Time {
	return u.updatedAt
}

func (u *User) IsActive() bool {
	return u.isActive
}

func (u *User) Deactivate() {
	u.isActive = false
	u.updatedAt = time.Now().UTC()
}

func (u *User) Activate() {
	u.isActive = true
	u.updatedAt = time.Now().UTC()
}

func (u *User) UpdateEmail(newEmail valueobject.Email) error {
	if newEmail.IsEmpty() {
		return errors.New("email cannot be empty")
	}

	u.email = newEmail
	u.updatedAt = time.Now().UTC()
	return nil
}

func (u *User) UpdatePassword(newPassword valueobject.Password) error {
	if newPassword.Hash() == "" {
		return errors.New("password cannot be empty")
	}

	u.password = newPassword
	u.updatedAt = time.Now().UTC()
	return nil
}

func (u *User) CanLogin() bool {
	return u.isActive
}

// Validate ensures user is in a valid state
// WHY: Defensive programming - catch invariant violations
func (u *User) Validate() error {
	if u.id.IsEmpty() {
		return errors.New("user ID cannot be empty")
	}

	if u.email.IsEmpty() {
		return errors.New("email cannot be empty")
	}

	if u.password.Hash() == "" {
		return errors.New("password cannot be empty")
	}

	return nil
}

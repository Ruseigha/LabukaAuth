package dto

import (
	"errors"
	"strings"
)

type SignupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *SignupRequest) Validate() error {
	// Trim whitespace
	r.Email = strings.TrimSpace(r.Email)
	r.Password = strings.TrimSpace(r.Password)

	if r.Email == "" {
		return errors.New("email is required")
	}

	if r.Password == "" {
		return errors.New("password is required")
	}

	return nil
}

// LoginRequest represents login HTTP request
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Validate validates login request
func (r *LoginRequest) Validate() error {
	r.Email = strings.TrimSpace(r.Email)
	r.Password = strings.TrimSpace(r.Password)

	if r.Email == "" {
		return errors.New("email is required")
	}

	if r.Password == "" {
		return errors.New("password is required")
	}

	return nil
}

// RefreshTokenRequest represents token refresh HTTP request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// Validate validates refresh token request
func (r *RefreshTokenRequest) Validate() error {
	r.RefreshToken = strings.TrimSpace(r.RefreshToken)

	if r.RefreshToken == "" {
		return errors.New("refresh_token is required")
	}

	return nil
}

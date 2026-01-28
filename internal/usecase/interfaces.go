package usecase

import "context"

type AuthUseCase interface {
	Signup(ctx context.Context, req SignupRequest) (*SignupResponse, error)
	Login(ctx context.Context, req LoginRequest) (*LoginResponse, error)
	ValidateToken(ctx context.Context, token string) (*TokenClaims, error)
	RefreshToken(ctx context.Context, refreshToken string) (*RefreshResponse, error)
}

// SignupRequest contains signup data
type SignupRequest struct {
	Email    string
	Password string
}

// SignupResponse contains signup result
type SignupResponse struct {
	UserID       string
	Email        string
	AccessToken  string
	RefreshToken string
}

// LoginRequest contains login credentials
type LoginRequest struct {
	Email    string
	Password string
}

// LoginResponse contains login result with tokens
type LoginResponse struct {
	UserID       string
	Email        string
	AccessToken  string
	RefreshToken string
}

// TokenClaims contains validated token data
type TokenClaims struct {
	UserID string
	Email  string
}

// RefreshResponse contains new tokens
type RefreshResponse struct {
	AccessToken  string
	RefreshToken string
}

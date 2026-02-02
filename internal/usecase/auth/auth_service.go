package auth

import (
	"context"

	"github.com/Ruseigha/LabukaAuth/internal/infrastructure/security"
	"github.com/Ruseigha/LabukaAuth/internal/repository"
	"github.com/Ruseigha/LabukaAuth/internal/usecase"
)

// AuthService aggregates all auth use cases
type AuthService struct {
	signupUC        *SignupUseCase
	loginUC         *LoginUseCase
	validateTokenUC *ValidateTokenUseCase
	refreshTokenUC  *RefreshTokenUseCase
}

// NewAuthService creates auth service with all use cases
func NewAuthService(
	userRepo repository.UserRepository,
	passwordHasher security.PasswordHasher,
	jwtGenerator security.JWTGenerator,
) *AuthService {
	return &AuthService{
		signupUC:        NewSignupUseCase(userRepo, passwordHasher, jwtGenerator),
		loginUC:         NewLoginUseCase(userRepo, passwordHasher, jwtGenerator),
		validateTokenUC: NewValidateTokenUseCase(jwtGenerator, userRepo),
		refreshTokenUC:  NewRefreshTokenUseCase(userRepo, jwtGenerator),
	}
}

// Signup registers a new user
func (s *AuthService) Signup(ctx context.Context, req usecase.SignupRequest) (*usecase.SignupResponse, error) {
	return s.signupUC.Execute(ctx, req)
}

// Login authenticates a user
func (s *AuthService) Login(ctx context.Context, req usecase.LoginRequest) (*usecase.LoginResponse, error) {
	return s.loginUC.Execute(ctx, req)
}

// ValidateToken validates an access token
func (s *AuthService) ValidateToken(ctx context.Context, token string) (*usecase.TokenClaims, error) {
	return s.validateTokenUC.Execute(ctx, token)
}

// RefreshToken generates new tokens from refresh token
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*usecase.RefreshResponse, error) {
	return s.refreshTokenUC.Execute(ctx, refreshToken)
}

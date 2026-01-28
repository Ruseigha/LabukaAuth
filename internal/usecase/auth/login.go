package auth

import (
	"context"
	"errors"
	"fmt"

	domainErrors "github.com/Ruseigha/LabukaAuth/internal/domain/errors"
	"github.com/Ruseigha/LabukaAuth/internal/domain/valueobject"
	"github.com/Ruseigha/LabukaAuth/internal/infrastructure/security"
	"github.com/Ruseigha/LabukaAuth/internal/repository"
	"github.com/Ruseigha/LabukaAuth/internal/usecase"
)

// LoginUseCase implements user authentication
type LoginUseCase struct {
	userRepo       repository.UserRepository
	passwordHasher security.PasswordHasher
	jwtGenerator   security.JWTGenerator
}

// NewLoginUseCase creates a new login use case
func NewLoginUseCase(
	userRepo repository.UserRepository,
	passwordHasher security.PasswordHasher,
	jwtGenerator security.JWTGenerator,
) *LoginUseCase {
	return &LoginUseCase{
		userRepo:       userRepo,
		passwordHasher: passwordHasher,
		jwtGenerator:   jwtGenerator,
	}
}

// Execute performs user login
func (uc *LoginUseCase) Execute(
	ctx context.Context,
	req usecase.LoginRequest,
) (*usecase.LoginResponse, error) {
	// Step 1: Validate email format
	email, err := valueobject.NewEmail(req.Email)
	if err != nil {
		// SECURITY: Don't reveal if email exists
		// WHY: Prevents email enumeration attacks
		return nil, domainErrors.NewUnauthorizedError("invalid credentials")
	}

	// Step 2: Find user by email
	user, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			// SECURITY: Same error as wrong password
			return nil, domainErrors.NewUnauthorizedError("invalid credentials")
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	// Step 3: Check if user is active
	// WHY: Business rule - inactive users can't login
	if !user.CanLogin() {
		return nil, domainErrors.NewForbiddenError("account is inactive")
	}

	// Step 4: Verify password
	// WHY: Core authentication - does password match?
	err = uc.passwordHasher.Compare(user.Password().Hash(), req.Password)
	if err != nil {
		// SECURITY: Same error as user not found
		return nil, domainErrors.NewUnauthorizedError("invalid credentials")
	}

	// Step 5: Generate tokens
	accessToken, err := uc.jwtGenerator.GenerateAccessToken(user.ID(), user.Email())
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := uc.jwtGenerator.GenerateRefreshToken(user.ID())
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Step 6: Return response
	return &usecase.LoginResponse{
		UserID:       user.ID().String(),
		Email:        user.Email().String(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

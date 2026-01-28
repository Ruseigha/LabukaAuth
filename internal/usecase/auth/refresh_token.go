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

// RefreshTokenUseCase implements token refresh logic
type RefreshTokenUseCase struct {
	userRepo     repository.UserRepository
	jwtGenerator security.JWTGenerator
}

// NewRefreshTokenUseCase creates a new refresh token use case
func NewRefreshTokenUseCase(
	userRepo repository.UserRepository,
	jwtGenerator security.JWTGenerator,
) *RefreshTokenUseCase {
	return &RefreshTokenUseCase{
		userRepo:     userRepo,
		jwtGenerator: jwtGenerator,
	}
}

// Execute refreshes access and refresh tokens
func (uc *RefreshTokenUseCase) Execute(
	ctx context.Context,
	refreshToken string,
) (*usecase.RefreshResponse, error) {
	// Step 1: Validate refresh token
	claims, err := uc.jwtGenerator.ValidateToken(refreshToken)
	if err != nil {
		return nil, domainErrors.NewUnauthorizedError("invalid refresh token")
	}

	// Step 2: Parse user ID
	userID, err := valueobject.NewUserIDFromString(claims.UserID)
	if err != nil {
		return nil, domainErrors.NewUnauthorizedError("invalid user ID in token")
	}

	// Step 3: Verify user exists
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, domainErrors.NewUnauthorizedError("user not found")
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	// Step 4: Check if user is active
	if !user.CanLogin() {
		return nil, domainErrors.NewForbiddenError("account is inactive")
	}

	// Step 5: Generate new access token
	newAccessToken, err := uc.jwtGenerator.GenerateAccessToken(user.ID(), user.Email())
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Step 6: Generate new refresh token
	newRefreshToken, err := uc.jwtGenerator.GenerateRefreshToken(user.ID())
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Step 7: Return new tokens
	return &usecase.RefreshResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

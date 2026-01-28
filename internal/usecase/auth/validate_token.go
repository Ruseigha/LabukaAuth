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

// ValidateTokenUseCase implements token validation
// WHY: Other services need to verify JWTs
type ValidateTokenUseCase struct {
	jwtGenerator security.JWTGenerator
	userRepo     repository.UserRepository
}

// NewValidateTokenUseCase creates a new validate token use case
func NewValidateTokenUseCase(
	jwtGenerator security.JWTGenerator,
	userRepo repository.UserRepository,
) *ValidateTokenUseCase {
	return &ValidateTokenUseCase{
		jwtGenerator: jwtGenerator,
		userRepo:     userRepo,
	}
}

// Execute validates a JWT token
func (uc *ValidateTokenUseCase) Execute(
	ctx context.Context,
	tokenString string,
) (*usecase.TokenClaims, error) {
	// Step 1: Validate JWT signature and claims
	claims, err := uc.jwtGenerator.ValidateToken(tokenString)
	if err != nil {
		return nil, domainErrors.NewUnauthorizedError("invalid token")
	}

	// Step 2: Parse user ID
	userID, err := valueobject.NewUserIDFromString(claims.UserID)
	if err != nil {
		return nil, domainErrors.NewUnauthorizedError("invalid user ID in token")
	}

	// Step 3: Verify user still exists and is active
	// WHY: User might be deleted or deactivated after token issued
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, domainErrors.NewUnauthorizedError("user not found")
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	// Step 4: Check if user can still login
	if !user.CanLogin() {
		return nil, domainErrors.NewForbiddenError("account is inactive")
	}

	// Step 5: Return validated claims
	return &usecase.TokenClaims{
		UserID: claims.UserID,
		Email:  claims.Email,
	}, nil
}

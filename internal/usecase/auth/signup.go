package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/Ruseigha/LabukaAuth/internal/domain/entity"
	domainErrors "github.com/Ruseigha/LabukaAuth/internal/domain/errors"
	"github.com/Ruseigha/LabukaAuth/internal/domain/valueobject"
	"github.com/Ruseigha/LabukaAuth/internal/infrastructure/security"
	"github.com/Ruseigha/LabukaAuth/internal/repository"
	"github.com/Ruseigha/LabukaAuth/internal/usecase"
)

// SignupUseCase implements user registration
// WHY: Orchestrates signup flow (validation, hashing, storage, token generation)
type SignupUseCase struct {
	userRepo       repository.UserRepository
	passwordHasher security.PasswordHasher
	jwtGenerator   security.JWTGenerator
}

// NewSignupUseCase creates a new signup use case
// WHY: Dependency injection - all dependencies passed in
func NewSignupUseCase(
	userRepo repository.UserRepository,
	passwordHasher security.PasswordHasher,
	jwtGenerator security.JWTGenerator,
) *SignupUseCase {
	return &SignupUseCase{
		userRepo:       userRepo,
		passwordHasher: passwordHasher,
		jwtGenerator:   jwtGenerator,
	}
}

// Execute performs user signup
// WHY: Single responsibility - this is the ONLY way to signup
func (uc *SignupUseCase) Execute(
	ctx context.Context,
	req usecase.SignupRequest,
) (*usecase.SignupResponse, error) {
	// Step 1: Validate input (create value objects)
	// WHY: Value objects enforce business rules
	email, err := valueobject.NewEmail(req.Email)
	if err != nil {
		return nil, domainErrors.NewInvalidInputError(
			"invalid email format",
			"email",
		)
	}

	password, err := valueobject.NewPassword(req.Password)
	if err != nil {
		return nil, domainErrors.NewInvalidInputError(
			err.Error(),
			"password",
		)
	}

	// Step 2: Check if user already exists
	// WHY: Business rule - emails must be unique
	exists, err := uc.userRepo.ExistsByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to check email existence: %w", err)
	}

	if exists {
		return nil, domainErrors.NewConflictError("email already in use")
	}

	// Step 3: Hash password
	// WHY: Never store plain text passwords
	hashedPassword, err := uc.passwordHasher.Hash(password.Hash())
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Step 4: Create hashed password value object
	hashedPasswordVO := valueobject.NewPasswordFromHash(hashedPassword)

	// Step 5: Create user entity
	// WHY: Entity enforces business rules and generates ID
	user, err := entity.NewUser(email, hashedPasswordVO)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Step 6: Save to repository
	// WHY: Persist the user
	if err := uc.userRepo.Create(ctx, user); err != nil {
		// Translate repository errors to domain errors
		if errors.Is(err, repository.ErrUserAlreadyExists) {
			// Race condition: User created between check and insert
			return nil, domainErrors.NewConflictError("email already in use")
		}
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Step 7: Generate tokens
	// WHY: User can immediately use the service after signup
	accessToken, err := uc.jwtGenerator.GenerateAccessToken(user.ID(), user.Email())
	if err != nil {
		// User is created but token generation failed
		// Log this for monitoring
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := uc.jwtGenerator.GenerateRefreshToken(user.ID())
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Step 8: Return response
	return &usecase.SignupResponse{
		UserID:       user.ID().String(),
		Email:        user.Email().String(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

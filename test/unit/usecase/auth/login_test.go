package auth_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Ruseigha/LabukaAuth/internal/domain/entity"
	domainErrors "github.com/Ruseigha/LabukaAuth/internal/domain/errors"
	"github.com/Ruseigha/LabukaAuth/internal/domain/valueobject"
	"github.com/Ruseigha/LabukaAuth/internal/repository"
	"github.com/Ruseigha/LabukaAuth/internal/usecase"
	"github.com/Ruseigha/LabukaAuth/internal/usecase/auth"
	"github.com/Ruseigha/LabukaAuth/test/unit/usecase/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestLoginUseCase_Success tests successful user login
func TestLoginUseCase_Success(t *testing.T) {
	// Arrange
	email, _ := valueobject.NewEmail("user@example.com")

	// Create a user with hashed password
	hashedPassword := valueobject.NewPasswordFromHash("hashed_SecureP@ss123")
	user, _ := entity.NewUser(email, hashedPassword)

	mockRepo := &mocks.MockUserRepository{
		FindByEmailFunc: func(ctx context.Context, em valueobject.Email) (*entity.User, error) {
			if em.Equals(email) {
				return user, nil
			}
			return nil, repository.ErrUserNotFound
		},
	}

	mockHasher := &mocks.MockPasswordHasher{
		CompareFunc: func(hashedPassword, plainPassword string) error {
			// Simulate successful password comparison
			if hashedPassword == "hashed_SecureP@ss123" && plainPassword == "SecureP@ss123" {
				return nil
			}
			return mocks.ErrPasswordMismatch
		},
	}

	mockJWT := &mocks.MockJWTGenerator{} // Uses default behavior

	loginUC := auth.NewLoginUseCase(mockRepo, mockHasher, mockJWT)

	req := usecase.LoginRequest{
		Email:    "user@example.com",
		Password: "SecureP@ss123",
	}

	// Act
	resp, err := loginUC.Execute(context.Background(), req)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.UserID)
	assert.Equal(t, "user@example.com", resp.Email)
	assert.NotEmpty(t, resp.AccessToken)
	assert.NotEmpty(t, resp.RefreshToken)

	// Verify repository and services were called
	assert.Equal(t, 1, mockRepo.FindByEmailCalls)
	assert.Equal(t, 1, mockHasher.CompareCalls)
	assert.Equal(t, 1, mockJWT.GenerateAccessTokenCalls)
	assert.Equal(t, 1, mockJWT.GenerateRefreshTokenCalls)
}

// TestLoginUseCase_InvalidEmail tests login with invalid email format
func TestLoginUseCase_InvalidEmail(t *testing.T) {
	tests := []struct {
		name  string
		email string
	}{
		{
			name:  "empty email",
			email: "",
		},
		{
			name:  "invalid format",
			email: "notanemail",
		},
		{
			name:  "missing @",
			email: "userexample.com",
		},
		{
			name:  "missing domain",
			email: "user@",
		},
		{
			name:  "whitespace only",
			email: "   ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := &mocks.MockUserRepository{}
			mockHasher := &mocks.MockPasswordHasher{}
			mockJWT := &mocks.MockJWTGenerator{}

			loginUC := auth.NewLoginUseCase(mockRepo, mockHasher, mockJWT)

			req := usecase.LoginRequest{
				Email:    tt.email,
				Password: "SecureP@ss123",
			}

			// Act
			resp, err := loginUC.Execute(context.Background(), req)

			// Assert
			require.Error(t, err)
			assert.Nil(t, resp)
			assert.True(t, errors.Is(err, domainErrors.ErrUnauthorized))

			// SECURITY: Should return "invalid credentials" not "invalid email"
			// WHY: Prevent email enumeration attacks
			assert.Contains(t, err.Error(), "invalid credentials")

			// Repository should NOT be called (validation failed)
			assert.Equal(t, 0, mockRepo.FindByEmailCalls)
			assert.Equal(t, 0, mockHasher.CompareCalls)
		})
	}
}

// TestLoginUseCase_UserNotFound tests login with non-existent email
func TestLoginUseCase_UserNotFound(t *testing.T) {
	// Arrange
	mockRepo := &mocks.MockUserRepository{
		FindByEmailFunc: func(ctx context.Context, email valueobject.Email) (*entity.User, error) {
			return nil, repository.ErrUserNotFound
		},
	}

	mockHasher := &mocks.MockPasswordHasher{}
	mockJWT := &mocks.MockJWTGenerator{}

	loginUC := auth.NewLoginUseCase(mockRepo, mockHasher, mockJWT)

	req := usecase.LoginRequest{
		Email:    "nonexistent@example.com",
		Password: "SecureP@ss123",
	}

	// Act
	resp, err := loginUC.Execute(context.Background(), req)

	// Assert
	require.Error(t, err)
	assert.Nil(t, resp)
	assert.True(t, errors.Is(err, domainErrors.ErrUnauthorized))

	// SECURITY: Same error as wrong password (prevents email enumeration)
	assert.Contains(t, err.Error(), "invalid credentials")

	// Repository was called but user not found
	assert.Equal(t, 1, mockRepo.FindByEmailCalls)

	// Password comparison should NOT be called (user doesn't exist)
	assert.Equal(t, 0, mockHasher.CompareCalls)
}

// TestLoginUseCase_WrongPassword tests login with incorrect password
func TestLoginUseCase_WrongPassword(t *testing.T) {
	// Arrange
	email, _ := valueobject.NewEmail("user@example.com")
	hashedPassword := valueobject.NewPasswordFromHash("hashed_SecureP@ss123")
	user, _ := entity.NewUser(email, hashedPassword)

	mockRepo := &mocks.MockUserRepository{
		FindByEmailFunc: func(ctx context.Context, em valueobject.Email) (*entity.User, error) {
			return user, nil
		},
	}

	mockHasher := &mocks.MockPasswordHasher{
		CompareFunc: func(hashedPassword, plainPassword string) error {
			// Simulate password mismatch
			return mocks.ErrPasswordMismatch
		},
	}

	mockJWT := &mocks.MockJWTGenerator{}

	loginUC := auth.NewLoginUseCase(mockRepo, mockHasher, mockJWT)

	req := usecase.LoginRequest{
		Email:    "user@example.com",
		Password: "WrongPassword123!", // Wrong password
	}

	// Act
	resp, err := loginUC.Execute(context.Background(), req)

	// Assert
	require.Error(t, err)
	assert.Nil(t, resp)
	assert.True(t, errors.Is(err, domainErrors.ErrUnauthorized))
	assert.Contains(t, err.Error(), "invalid credentials")

	// User was found and password was compared
	assert.Equal(t, 1, mockRepo.FindByEmailCalls)
	assert.Equal(t, 1, mockHasher.CompareCalls)

	// Tokens should NOT be generated (authentication failed)
	assert.Equal(t, 0, mockJWT.GenerateAccessTokenCalls)
	assert.Equal(t, 0, mockJWT.GenerateRefreshTokenCalls)
}

// TestLoginUseCase_InactiveUser tests login with deactivated account
func TestLoginUseCase_InactiveUser(t *testing.T) {
	// Arrange
	email, _ := valueobject.NewEmail("inactive@example.com")
	hashedPassword := valueobject.NewPasswordFromHash("hashed_SecureP@ss123")

	// Create user then deactivate
	user, _ := entity.NewUser(email, hashedPassword)
	user.Deactivate() // User is now inactive

	mockRepo := &mocks.MockUserRepository{
		FindByEmailFunc: func(ctx context.Context, em valueobject.Email) (*entity.User, error) {
			return user, nil
		},
	}

	mockHasher := &mocks.MockPasswordHasher{
		CompareFunc: func(hashedPassword, plainPassword string) error {
			return nil // Password is correct
		},
	}

	mockJWT := &mocks.MockJWTGenerator{}

	loginUC := auth.NewLoginUseCase(mockRepo, mockHasher, mockJWT)

	req := usecase.LoginRequest{
		Email:    "inactive@example.com",
		Password: "SecureP@ss123",
	}

	// Act
	resp, err := loginUC.Execute(context.Background(), req)

	// Assert
	require.Error(t, err)
	assert.Nil(t, resp)
	assert.True(t, errors.Is(err, domainErrors.ErrForbidden))
	assert.Contains(t, err.Error(), "account is inactive")

	// User was found, but login prevented due to inactive status
	assert.Equal(t, 1, mockRepo.FindByEmailCalls)

	// Password comparison should NOT happen (user inactive check comes first)
	// NOTE: Depends on your implementation order
	// If you check password first, this would be 1
}

// TestLoginUseCase_CaseInsensitiveEmail tests email case insensitivity
func TestLoginUseCase_CaseInsensitiveEmail(t *testing.T) {
	// Arrange
	email, _ := valueobject.NewEmail("user@example.com") // Lowercase
	hashedPassword := valueobject.NewPasswordFromHash("hashed_SecureP@ss123")
	user, _ := entity.NewUser(email, hashedPassword)

	mockRepo := &mocks.MockUserRepository{
		FindByEmailFunc: func(ctx context.Context, em valueobject.Email) (*entity.User, error) {
			// Email value object normalizes to lowercase
			if em.String() == "user@example.com" {
				return user, nil
			}
			return nil, repository.ErrUserNotFound
		},
	}

	mockHasher := &mocks.MockPasswordHasher{
		CompareFunc: func(hashedPassword, plainPassword string) error {
			return nil
		},
	}

	mockJWT := &mocks.MockJWTGenerator{}

	loginUC := auth.NewLoginUseCase(mockRepo, mockHasher, mockJWT)

	// Test with different cases
	testCases := []string{
		"user@example.com",
		"User@Example.COM",
		"USER@EXAMPLE.COM",
		"uSeR@eXaMpLe.CoM",
	}

	for _, emailInput := range testCases {
		t.Run(emailInput, func(t *testing.T) {
			req := usecase.LoginRequest{
				Email:    emailInput,
				Password: "SecureP@ss123",
			}

			// Act
			resp, err := loginUC.Execute(context.Background(), req)

			// Assert
			require.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, "user@example.com", resp.Email) // Normalized
		})
	}
}

// TestLoginUseCase_RepositoryError tests database errors
func TestLoginUseCase_RepositoryError(t *testing.T) {
	// Arrange
	mockRepo := &mocks.MockUserRepository{
		FindByEmailFunc: func(ctx context.Context, email valueobject.Email) (*entity.User, error) {
			return nil, errors.New("database connection lost")
		},
	}

	mockHasher := &mocks.MockPasswordHasher{}
	mockJWT := &mocks.MockJWTGenerator{}

	loginUC := auth.NewLoginUseCase(mockRepo, mockHasher, mockJWT)

	req := usecase.LoginRequest{
		Email:    "user@example.com",
		Password: "SecureP@ss123",
	}

	// Act
	resp, err := loginUC.Execute(context.Background(), req)

	// Assert
	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "failed to find user")

	// Repository was called
	assert.Equal(t, 1, mockRepo.FindByEmailCalls)
}

// TestLoginUseCase_TokenGenerationError tests JWT generation failures
func TestLoginUseCase_TokenGenerationError(t *testing.T) {
	tests := []struct {
		name      string
		setupMock func() *mocks.MockJWTGenerator
	}{
		{
			name: "access token generation fails",
			setupMock: func() *mocks.MockJWTGenerator {
				return &mocks.MockJWTGenerator{
					GenerateAccessTokenFunc: func(userID valueobject.UserID, email valueobject.Email) (string, error) {
						return "", errors.New("signing key not found")
					},
				}
			},
		},
		{
			name: "refresh token generation fails",
			setupMock: func() *mocks.MockJWTGenerator {
				return &mocks.MockJWTGenerator{
					GenerateAccessTokenFunc: func(userID valueobject.UserID, email valueobject.Email) (string, error) {
						return "access_token", nil // Success
					},
					GenerateRefreshTokenFunc: func(userID valueobject.UserID) (string, error) {
						return "", errors.New("signing key not found")
					},
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			email, _ := valueobject.NewEmail("user@example.com")
			hashedPassword := valueobject.NewPasswordFromHash("hashed_SecureP@ss123")
			user, _ := entity.NewUser(email, hashedPassword)

			mockRepo := &mocks.MockUserRepository{
				FindByEmailFunc: func(ctx context.Context, em valueobject.Email) (*entity.User, error) {
					return user, nil
				},
			}

			mockHasher := &mocks.MockPasswordHasher{
				CompareFunc: func(hashedPassword, plainPassword string) error {
					return nil // Password correct
				},
			}

			mockJWT := tt.setupMock()

			loginUC := auth.NewLoginUseCase(mockRepo, mockHasher, mockJWT)

			req := usecase.LoginRequest{
				Email:    "user@example.com",
				Password: "SecureP@ss123",
			}

			// Act
			resp, err := loginUC.Execute(context.Background(), req)

			// Assert
			require.Error(t, err)
			assert.Nil(t, resp)
			assert.Contains(t, err.Error(), "failed to generate")

			// User was authenticated but token generation failed
			assert.Equal(t, 1, mockRepo.FindByEmailCalls)
			assert.Equal(t, 1, mockHasher.CompareCalls)
		})
	}
}

// TestLoginUseCase_PasswordComparisonTiming tests timing attack prevention
// NOTE: This is more of a demonstration - real timing attacks are hard to test
func TestLoginUseCase_PasswordComparisonTiming(t *testing.T) {
	// This test demonstrates that we use constant-time comparison
	// In practice, bcrypt.CompareHashAndPassword is constant-time

	email, _ := valueobject.NewEmail("user@example.com")
	hashedPassword := valueobject.NewPasswordFromHash("hashed_SecureP@ss123")
	user, _ := entity.NewUser(email, hashedPassword)

	mockRepo := &mocks.MockUserRepository{
		FindByEmailFunc: func(ctx context.Context, em valueobject.Email) (*entity.User, error) {
			return user, nil
		},
	}

	callCount := 0
	mockHasher := &mocks.MockPasswordHasher{
		CompareFunc: func(hashedPassword, plainPassword string) error {
			callCount++
			// Simulate constant-time comparison
			// In reality, bcrypt does this
			return mocks.ErrPasswordMismatch
		},
	}

	mockJWT := &mocks.MockJWTGenerator{}

	loginUC := auth.NewLoginUseCase(mockRepo, mockHasher, mockJWT)

	// Try multiple wrong passwords
	passwords := []string{
		"a",             // Very short
		"WrongPassword", // Wrong but similar length
		"VeryLongWrongPasswordThatIsDefinitelyNotCorrect", // Very long
	}

	for _, pwd := range passwords {
		req := usecase.LoginRequest{
			Email:    "user@example.com",
			Password: pwd,
		}

		loginUC.Execute(context.Background(), req)
	}

	// Verify password comparison happened for all attempts
	assert.Equal(t, len(passwords), callCount)

	// NOTE: In production, bcrypt.CompareHashAndPassword takes roughly
	// the same time regardless of password correctness
}

// TestLoginUseCase_MultipleFailedAttempts tests behavior with repeated failures
// NOTE: Rate limiting should be handled at middleware level
func TestLoginUseCase_MultipleFailedAttempts(t *testing.T) {
	// Arrange
	email, _ := valueobject.NewEmail("user@example.com")
	hashedPassword := valueobject.NewPasswordFromHash("hashed_SecureP@ss123")
	user, _ := entity.NewUser(email, hashedPassword)

	mockRepo := &mocks.MockUserRepository{
		FindByEmailFunc: func(ctx context.Context, em valueobject.Email) (*entity.User, error) {
			return user, nil
		},
	}

	mockHasher := &mocks.MockPasswordHasher{
		CompareFunc: func(hashedPassword, plainPassword string) error {
			return mocks.ErrPasswordMismatch // Always wrong
		},
	}

	mockJWT := &mocks.MockJWTGenerator{}

	loginUC := auth.NewLoginUseCase(mockRepo, mockHasher, mockJWT)

	// Act - Try multiple times
	for i := 0; i < 5; i++ {
		req := usecase.LoginRequest{
			Email:    "user@example.com",
			Password: "WrongPassword",
		}

		resp, err := loginUC.Execute(context.Background(), req)

		// Assert - Each attempt should fail consistently
		require.Error(t, err)
		assert.Nil(t, resp)
		assert.True(t, errors.Is(err, domainErrors.ErrUnauthorized))
	}

	// Verify all attempts called the repository
	assert.Equal(t, 5, mockRepo.FindByEmailCalls)
	assert.Equal(t, 5, mockHasher.CompareCalls)

	// NOTE: Account lockout after N failed attempts should be
	// implemented as a separate concern (middleware or use case decorator)
}

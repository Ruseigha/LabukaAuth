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

// TestSignupUseCase_Success tests successful user signup
func TestSignupUseCase_Success(t *testing.T) {
    // Arrange
    mockRepo := &mocks.MockUserRepository{
        ExistsByEmailFunc: func(ctx context.Context, email valueobject.Email) (bool, error) {
            return false, nil  // Email doesn't exist
        },
        CreateFunc: func(ctx context.Context, user *entity.User) error {
            return nil  // Success
        },
    }
    
    mockHasher := &mocks.MockPasswordHasher{}  // Uses default behavior
    mockJWT := &mocks.MockJWTGenerator{}       // Uses default behavior
    
    signupUC := auth.NewSignupUseCase(mockRepo, mockHasher, mockJWT)
    
    req := usecase.SignupRequest{
        Email:    "newuser@example.com",
        Password: "SecureP@ss123",
    }
    
    // Act
    resp, err := signupUC.Execute(context.Background(), req)
    
    // Assert
    require.NoError(t, err)
    assert.NotNil(t, resp)
    assert.NotEmpty(t, resp.UserID)
    assert.Equal(t, "newuser@example.com", resp.Email)
    assert.NotEmpty(t, resp.AccessToken)
    assert.NotEmpty(t, resp.RefreshToken)
    
    // Verify repository was called
    assert.Equal(t, 1, mockRepo.ExistsByEmailCalls)
    assert.Equal(t, 1, mockRepo.CreateCalls)
    assert.Equal(t, 1, mockHasher.HashCalls)
    assert.Equal(t, 1, mockJWT.GenerateAccessTokenCalls)
    assert.Equal(t, 1, mockJWT.GenerateRefreshTokenCalls)
}

// TestSignupUseCase_InvalidEmail tests signup with invalid email
func TestSignupUseCase_InvalidEmail(t *testing.T) {
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
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Arrange
            mockRepo := &mocks.MockUserRepository{}
            mockHasher := &mocks.MockPasswordHasher{}
            mockJWT := &mocks.MockJWTGenerator{}
            
            signupUC := auth.NewSignupUseCase(mockRepo, mockHasher, mockJWT)
            
            req := usecase.SignupRequest{
                Email:    tt.email,
                Password: "SecureP@ss123",
            }
            
            // Act
            resp, err := signupUC.Execute(context.Background(), req)
            
            // Assert
            require.Error(t, err)
            assert.Nil(t, resp)
            assert.True(t, errors.Is(err, domainErrors.ErrInvalidInput))
            
            // Repository should NOT be called
            assert.Equal(t, 0, mockRepo.ExistsByEmailCalls)
            assert.Equal(t, 0, mockRepo.CreateCalls)
        })
    }
}

// TestSignupUseCase_InvalidPassword tests signup with invalid password
func TestSignupUseCase_InvalidPassword(t *testing.T) {
    tests := []struct {
        name     string
        password string
    }{
        {
            name:     "too short",
            password: "Pass1!",
        },
        {
            name:     "no uppercase",
            password: "password123!",
        },
        {
            name:     "no lowercase",
            password: "PASSWORD123!",
        },
        {
            name:     "no number",
            password: "Password!",
        },
        {
            name:     "no special char",
            password: "Password123",
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Arrange
            mockRepo := &mocks.MockUserRepository{}
            mockHasher := &mocks.MockPasswordHasher{}
            mockJWT := &mocks.MockJWTGenerator{}
            
            signupUC := auth.NewSignupUseCase(mockRepo, mockHasher, mockJWT)
            
            req := usecase.SignupRequest{
                Email:    "user@example.com",
                Password: tt.password,
            }
            
            // Act
            resp, err := signupUC.Execute(context.Background(), req)
            
            // Assert
            require.Error(t, err)
            assert.Nil(t, resp)
            assert.True(t, errors.Is(err, domainErrors.ErrInvalidInput))
            
            // Repository should NOT be called
            assert.Equal(t, 0, mockRepo.ExistsByEmailCalls)
        })
    }
}

// TestSignupUseCase_EmailAlreadyExists tests duplicate email
func TestSignupUseCase_EmailAlreadyExists(t *testing.T) {
    // Arrange
    mockRepo := &mocks.MockUserRepository{
        ExistsByEmailFunc: func(ctx context.Context, email valueobject.Email) (bool, error) {
            return true, nil  // Email exists!
        },
    }
    mockHasher := &mocks.MockPasswordHasher{}
    mockJWT := &mocks.MockJWTGenerator{}
    
    signupUC := auth.NewSignupUseCase(mockRepo, mockHasher, mockJWT)
    
    req := usecase.SignupRequest{
        Email:    "existing@example.com",
        Password: "SecureP@ss123",
    }
    
    // Act
    resp, err := signupUC.Execute(context.Background(), req)
    
    // Assert
    require.Error(t, err)
    assert.Nil(t, resp)
    assert.True(t, errors.Is(err, domainErrors.ErrConflict))
    
    // Should check existence but NOT create
    assert.Equal(t, 1, mockRepo.ExistsByEmailCalls)
    assert.Equal(t, 0, mockRepo.CreateCalls)
}

// TestSignupUseCase_RepositoryCreateError tests database errors
func TestSignupUseCase_RepositoryCreateError(t *testing.T) {
    // Arrange
    mockRepo := &mocks.MockUserRepository{
        ExistsByEmailFunc: func(ctx context.Context, email valueobject.Email) (bool, error) {
            return false, nil
        },
        CreateFunc: func(ctx context.Context, user *entity.User) error {
            return errors.New("database connection lost")
        },
    }
    mockHasher := &mocks.MockPasswordHasher{}
    mockJWT := &mocks.MockJWTGenerator{}
    
    signupUC := auth.NewSignupUseCase(mockRepo, mockHasher, mockJWT)
    
    req := usecase.SignupRequest{
        Email:    "user@example.com",
        Password: "SecureP@ss123",
    }
    
    // Act
    resp, err := signupUC.Execute(context.Background(), req)
    
    // Assert
    require.Error(t, err)
    assert.Nil(t, resp)
    assert.Contains(t, err.Error(), "failed to create user")
}

// TestSignupUseCase_RaceCondition tests duplicate detection race condition
func TestSignupUseCase_RaceCondition(t *testing.T) {
    // Arrange - Simulates race condition
    // ExistsByEmail returns false, but Create fails with duplicate
    mockRepo := &mocks.MockUserRepository{
        ExistsByEmailFunc: func(ctx context.Context, email valueobject.Email) (bool, error) {
            return false, nil  // Check passes
        },
        CreateFunc: func(ctx context.Context, user *entity.User) error {
            // But insert fails (another thread created user)
            return repository.ErrUserAlreadyExists
        },
    }
    mockHasher := &mocks.MockPasswordHasher{}
    mockJWT := &mocks.MockJWTGenerator{}
    
    signupUC := auth.NewSignupUseCase(mockRepo, mockHasher, mockJWT)
    
    req := usecase.SignupRequest{
        Email:    "user@example.com",
        Password: "SecureP@ss123",
    }
    
    // Act
    resp, err := signupUC.Execute(context.Background(), req)
    
    // Assert
    require.Error(t, err)
    assert.Nil(t, resp)
    assert.True(t, errors.Is(err, domainErrors.ErrConflict))
}

// TestSignupUseCase_PasswordHashingError tests password hashing failure
func TestSignupUseCase_PasswordHashingError(t *testing.T) {
    // Arrange
    mockRepo := &mocks.MockUserRepository{
        ExistsByEmailFunc: func(ctx context.Context, email valueobject.Email) (bool, error) {
            return false, nil
        },
    }
    mockHasher := &mocks.MockPasswordHasher{
        HashFunc: func(password string) (string, error) {
            return "", errors.New("hashing algorithm failure")
        },
    }
    mockJWT := &mocks.MockJWTGenerator{}
    
    signupUC := auth.NewSignupUseCase(mockRepo, mockHasher, mockJWT)
    
    req := usecase.SignupRequest{
        Email:    "user@example.com",
        Password: "SecureP@ss123",
    }
    
    // Act
    resp, err := signupUC.Execute(context.Background(), req)
    
    // Assert
    require.Error(t, err)
    assert.Nil(t, resp)
    assert.Contains(t, err.Error(), "failed to hash password")
    
    // Should NOT create user if hashing fails
    assert.Equal(t, 0, mockRepo.CreateCalls)
}

// TestSignupUseCase_TokenGenerationError tests JWT generation failure
func TestSignupUseCase_TokenGenerationError(t *testing.T) {
    // Arrange
    mockRepo := &mocks.MockUserRepository{
        ExistsByEmailFunc: func(ctx context.Context, email valueobject.Email) (bool, error) {
            return false, nil
        },
        CreateFunc: func(ctx context.Context, user *entity.User) error {
            return nil  // User created successfully
        },
    }
    mockHasher := &mocks.MockPasswordHasher{}
    mockJWT := &mocks.MockJWTGenerator{
        GenerateAccessTokenFunc: func(userID valueobject.UserID, email valueobject.Email) (string, error) {
            return "", errors.New("JWT signing key not found")
        },
    }
    
    signupUC := auth.NewSignupUseCase(mockRepo, mockHasher, mockJWT)
    
    req := usecase.SignupRequest{
        Email:    "user@example.com",
        Password: "SecureP@ss123",
    }
    
    // Act
    resp, err := signupUC.Execute(context.Background(), req)
    
    // Assert
    require.Error(t, err)
    assert.Nil(t, resp)
    assert.Contains(t, err.Error(), "failed to generate access token")
    
    // User WAS created (token generation happens after)
    assert.Equal(t, 1, mockRepo.CreateCalls)
}
package mocks

import (
	"errors"

	"github.com/Ruseigha/LabukaAuth/internal/domain/valueobject"
	"github.com/Ruseigha/LabukaAuth/internal/infrastructure/security"
)

// MockJWTGenerator is a mock implementation of JWTGenerator
type MockJWTGenerator struct {
	GenerateAccessTokenFunc  func(userID valueobject.UserID, email valueobject.Email) (string, error)
	GenerateRefreshTokenFunc func(userID valueobject.UserID) (string, error)
	ValidateTokenFunc        func(tokenString string) (*security.Claims, error)

	GenerateAccessTokenCalls  int
	GenerateRefreshTokenCalls int
	ValidateTokenCalls        int
}

// GenerateAccessToken implements security.JWTGenerator
func (m *MockJWTGenerator) GenerateAccessToken(
	userID valueobject.UserID,
	email valueobject.Email,
) (string, error) {
	m.GenerateAccessTokenCalls++
	if m.GenerateAccessTokenFunc != nil {
		return m.GenerateAccessTokenFunc(userID, email)
	}
	// Default: return predictable token
	return "access_token_" + userID.String(), nil
}

// GenerateRefreshToken implements security.JWTGenerator
func (m *MockJWTGenerator) GenerateRefreshToken(userID valueobject.UserID) (string, error) {
	m.GenerateRefreshTokenCalls++
	if m.GenerateRefreshTokenFunc != nil {
		return m.GenerateRefreshTokenFunc(userID)
	}
	// Default: return predictable token
	return "refresh_token_" + userID.String(), nil
}

// ValidateToken implements security.JWTGenerator
func (m *MockJWTGenerator) ValidateToken(tokenString string) (*security.Claims, error) {
	m.ValidateTokenCalls++
	if m.ValidateTokenFunc != nil {
		return m.ValidateTokenFunc(tokenString)
	}
	return nil, errors.New("invalid token")
}

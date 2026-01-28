package security

import (
	"errors"
	"time"

	"github.com/Ruseigha/LabukaAuth/internal/domain/valueobject"
	"github.com/golang-jwt/jwt/v5"
)

type JWTGenerator interface {
	GenerateAccessToken(userID valueobject.UserID, email valueobject.Email) (string, error)
	GenerateRefreshToken(userID valueobject.UserID) (string, error)
	ValidateToken(tokenString string) (*Claims, error)
}

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

type JWTGeneratorImpl struct {
	secretKey          []byte        // Signing key
	accessTokenExpiry  time.Duration // Access token lifetime
	refreshTokenExpiry time.Duration // Refresh token lifetime
	issuer             string        // JWT issuer
}

// NewJWTGenerator creates a new JWT generator
func NewJWTGenerator(
	secretKey string,
	accessTokenExpiry time.Duration,
	refreshTokenExpiry time.Duration,
	issuer string,
) *JWTGeneratorImpl {
	return &JWTGeneratorImpl{
		secretKey:          []byte(secretKey),
		accessTokenExpiry:  accessTokenExpiry,
		refreshTokenExpiry: refreshTokenExpiry,
		issuer:             issuer,
	}
}

// GenerateAccessToken creates an access token
func (g *JWTGeneratorImpl) GenerateAccessToken(
	userID valueobject.UserID,
	email valueobject.Email,
) (string, error) {
	now := time.Now()
	expiresAt := now.Add(g.accessTokenExpiry)

	// Create claims
	claims := Claims{
		UserID: userID.String(),
		Email:  email.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    g.issuer,
			Subject:   userID.String(),
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token
	tokenString, err := token.SignedString(g.secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// GenerateRefreshToken creates a refresh token
// WHY: Refresh tokens don't need email (less data in token)
func (g *JWTGeneratorImpl) GenerateRefreshToken(
	userID valueobject.UserID,
) (string, error) {
	now := time.Now()
	expiresAt := now.Add(g.refreshTokenExpiry)

	claims := Claims{
		UserID: userID.String(),
		// No email in refresh token
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    g.issuer,
			Subject:   userID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(g.secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates JWT and returns claims
func (g *JWTGeneratorImpl) ValidateToken(tokenString string) (*Claims, error) {
	// Parse token
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			// Verify signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return g.secretKey, nil
		},
	)

	if err != nil {
		return nil, err
	}

	// Extract claims
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
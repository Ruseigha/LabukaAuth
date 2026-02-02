package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/Ruseigha/LabukaAuth/internal/usecase"
)

// contextKey is a custom type for context keys
// WHY: Avoid collisions with other packages
type contextKey string

const (
	// UserIDKey is the context key for user ID
	UserIDKey contextKey = "user_id"

	// EmailKey is the context key for email
	EmailKey contextKey = "email"
)

// Auth validates JWT tokens
// WHY: Protect endpoints that require authentication
func Auth(authService usecase.AuthUseCase) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract token from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				respondUnauthorized(w, "missing authorization header")
				return
			}

			// Check Bearer prefix
			const bearerPrefix = "Bearer "
			if !strings.HasPrefix(authHeader, bearerPrefix) {
				respondUnauthorized(w, "invalid authorization header format")
				return
			}

			// Extract token
			token := authHeader[len(bearerPrefix):]
			if token == "" {
				respondUnauthorized(w, "missing token")
				return
			}

			// Validate token
			claims, err := authService.ValidateToken(r.Context(), token)
			if err != nil {
				respondUnauthorized(w, "invalid token")
				return
			}

			// Add user info to context
			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, EmailKey, claims.Email)

			// Call next handler with enriched context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// respondUnauthorized sends 401 response
func respondUnauthorized(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(`{"error":"unauthorized","message":"` + message + `","code":401}`))
}

// GetUserIDFromContext extracts user ID from context
// WHY: Helper for handlers to get authenticated user ID
func GetUserIDFromContext(ctx context.Context) string {
	userID, _ := ctx.Value(UserIDKey).(string)
	return userID
}

// GetEmailFromContext extracts email from context
func GetEmailFromContext(ctx context.Context) string {
	email, _ := ctx.Value(EmailKey).(string)
	return email
}

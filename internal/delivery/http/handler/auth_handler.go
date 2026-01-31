package handler

import (
	"encoding/json"
	"errors"

	"net/http"

	"github.com/Ruseigha/LabukaAuth/internal/delivery/http/dto"
	"github.com/Ruseigha/LabukaAuth/internal/usecase"
)

type AuthHandler struct {
	authService usecase.AuthUseCase
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Helper: respondJSON sends JSON response
func respondJSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// Helper: respondError sends error response
func respondError(w http.ResponseWriter, statusCode int, message string, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	errResp := dto.NewErrorResponse(err, message, statusCode)
	json.NewEncoder(w).Encode(errResp)
}

// Helper: extractBearerToken extracts JWT from Authorization header
func extractBearerToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	// Expected format: "Bearer <token>"
	const bearerPrefix = "Bearer "
	if len(authHeader) < len(bearerPrefix) {
		return ""
	}

	if authHeader[:len(bearerPrefix)] != bearerPrefix {
		return ""
	}

	return authHeader[len(bearerPrefix):]
}

func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	// Step 1: Parse request body
	var req dto.SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	// Step 2: Validate HTTP request
	if err := req.Validate(); err != nil {
		respondError(w, http.StatusBadRequest, "validation failed", err)
		return
	}

	// Step 3: Call use case
	resp, err := h.authService.Signup(r.Context(), usecase.SignupRequest{
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		// Map domain error to HTTP status
		statusCode := dto.MapDomainErrorToHTTP(err)
		respondError(w, statusCode, "signup failed", err)
		return
	}

	// Step 4: Return success response
	respondJSON(w, http.StatusCreated, dto.AuthResponse{
		UserID:       resp.UserID,
		Email:        resp.Email,
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// Parse request
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	// Validate
	if err := req.Validate(); err != nil {
		respondError(w, http.StatusBadRequest, "validation failed", err)
		return
	}

	// Call use case
	resp, err := h.authService.Login(r.Context(), usecase.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		statusCode := dto.MapDomainErrorToHTTP(err)
		respondError(w, statusCode, "login failed", err)
		return
	}

	// Return response
	respondJSON(w, http.StatusOK, dto.AuthResponse{
		UserID:       resp.UserID,
		Email:        resp.Email,
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
	})
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	// Parse request
	var req dto.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	// Validate
	if err := req.Validate(); err != nil {
		respondError(w, http.StatusBadRequest, "validation failed", err)
		return
	}

	// Call use case
	resp, err := h.authService.RefreshToken(r.Context(), req.RefreshToken)

	if err != nil {
		statusCode := dto.MapDomainErrorToHTTP(err)
		respondError(w, statusCode, "token refresh failed", err)
		return
	}

	// Return new tokens (no email in refresh response)
	respondJSON(w, http.StatusOK, dto.AuthResponse{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
	})
}

func (h *AuthHandler) ValidateToken(w http.ResponseWriter, r *http.Request) {
	// Extract token from Authorization header
	token := extractBearerToken(r)
	if token == "" {
		respondError(w, http.StatusUnauthorized, "missing authorization header", errors.New("no token provided"))
		return
	}

	// Call use case
	claims, err := h.authService.ValidateToken(r.Context(), token)

	if err != nil {
		statusCode := dto.MapDomainErrorToHTTP(err)
		respondError(w, statusCode, "token validation failed", err)
		return
	}

	// Return validation result
	respondJSON(w, http.StatusOK, dto.ValidateTokenResponse{
		Valid:  true,
		UserID: claims.UserID,
		Email:  claims.Email,
	})
}

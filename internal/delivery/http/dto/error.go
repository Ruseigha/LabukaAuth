package dto

import (
	"errors"
	"net/http"

	domainerrors "github.com/Ruseigha/LabukaAuth/internal/domain/errors"
)

// ErrorResponse represents HTTP error response
// WHY: Consistent error format across all endpoints
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// NewErrorResponse creates error response
func NewErrorResponse(err error, message string, code int) *ErrorResponse {
	return &ErrorResponse{
		Error:   err.Error(),
		Message: message,
		Code:    code,
	}
}

// MapDomainErrorToHTTP maps domain errors to HTTP status codes
// WHY: Translate business errors to appropriate HTTP codes
func MapDomainErrorToHTTP(err error) int {
	// Check domain error types
	if errors.Is(err, domainerrors.ErrInvalidInput) {
		return http.StatusBadRequest // 400
	}

	if errors.Is(err, domainerrors.ErrUnauthorized) {
		return http.StatusUnauthorized // 401
	}

	if errors.Is(err, domainerrors.ErrForbidden) {
		return http.StatusForbidden // 403
	}

	if errors.Is(err, domainerrors.ErrNotFound) {
		return http.StatusNotFound // 404
	}

	if errors.Is(err, domainerrors.ErrConflict) {
		return http.StatusConflict // 409
	}

	// Default to internal server error
	return http.StatusInternalServerError // 500
}

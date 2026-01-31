package handler

import (
	"net/http"

	"github.com/Ruseigha/LabukaAuth/internal/delivery/http/dto"
)

// HealthHandler handles health check requests
type HealthHandler struct {
	version string
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(version string) *HealthHandler {
	return &HealthHandler{
		version: version,
	}
}

func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, dto.HealthResponse{
		Status:  "healthy",
		Service: "auth-service",
		Version: h.version,
	})
}

func (h *HealthHandler) Ready(w http.ResponseWriter, r *http.Request) {
	// TODO: Check database connection, etc.
	respondJSON(w, http.StatusOK, dto.HealthResponse{
		Status:  "ready",
		Service: "auth-service",
		Version: h.version,
	})
}

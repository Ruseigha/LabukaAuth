package http

import (
	"net/http"

	"github.com/Ruseigha/LabukaAuth/internal/delivery/http/handler"
	"github.com/Ruseigha/LabukaAuth/internal/delivery/http/middleware"
	"github.com/Ruseigha/LabukaAuth/internal/usecase"
	"github.com/gorilla/mux"
)

// SetupRouter creates and configures the HTTP router
func SetupRouter(authService usecase.AuthUseCase, version string) http.Handler {
	// Create router
	r := mux.NewRouter()

	// Create handlers
	authHandler := handler.NewAuthHandler(authService)
	healthHandler := handler.NewHealthHandler(version)

	// Health check routes (no auth required)
	r.HandleFunc("/health", healthHandler.Health).Methods(http.MethodGet)
	r.HandleFunc("/ready", healthHandler.Ready).Methods(http.MethodGet)

	// API v1 routes
	api := r.PathPrefix("/api/v1").Subrouter()

	// Auth routes (public - no auth required)
	api.HandleFunc("/auth/signup", authHandler.Signup).Methods(http.MethodPost)
	api.HandleFunc("/auth/login", authHandler.Login).Methods(http.MethodPost)
	api.HandleFunc("/auth/refresh", authHandler.RefreshToken).Methods(http.MethodPost)

	// Protected routes (require authentication)
	protected := api.PathPrefix("").Subrouter()
	protected.Use(middleware.Auth(authService)) // Apply auth middleware
	protected.HandleFunc("/auth/validate", authHandler.ValidateToken).Methods(http.MethodGet)

	// Apply global middleware (in order)
	handler := middleware.Recovery(r)                // Outermost: catch panics
	handler = middleware.Logger(handler)             // Log all requests
	handler = middleware.CORS(handler)               // Add CORS headers
	handler = middleware.RateLimit(100, 20)(handler) // 100 req/min, burst 20

	return handler
}

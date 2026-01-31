package dto

// AuthResponse represents authentication response with tokens
// WHY: Consistent response structure for signup/login/refresh
type AuthResponse struct {
	UserID       string `json:"user_id"`
	Email        string `json:"email,omitempty"` // Omit in refresh response
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// ValidateTokenResponse represents token validation response
type ValidateTokenResponse struct {
	Valid  bool   `json:"valid"`
	UserID string `json:"user_id,omitempty"`
	Email  string `json:"email,omitempty"`
}

// HealthResponse represents health check response
type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
	Version string `json:"version"`
}

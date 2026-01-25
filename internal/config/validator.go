package config

import (
	"errors"
	"fmt"
	"strings"
)

// Validate checks if configuration is valid
// WHY: Fail fast on startup, not at runtime
func Validate(cfg *Config) error {
	var errs []error

	// Validate App config
	if err := validateApp(&cfg.App); err != nil {
		errs = append(errs, err)
	}

	// Validate Server config
	if err := validateServer(&cfg.Server); err != nil {
		errs = append(errs, err)
	}

	// Validate Database config
	if err := validateDatabase(&cfg.Database); err != nil {
		errs = append(errs, err)
	}

	// Validate JWT config
	if err := validateJWT(&cfg.JWT); err != nil {
		errs = append(errs, err)
	}

	// Validate Logger config
	if err := validateLogger(&cfg.Logger); err != nil {
		errs = append(errs, err)
	}

	// Combine all errors
	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

// validateApp validates application configuration
func validateApp(cfg *AppConfig) error {
	var errs []error

	if cfg.Name == "" {
		errs = append(errs, errors.New("app name is required"))
	}

	// Validate environment
	validEnvs := []string{"development", "dev", "staging", "production", "prod", "test"}
	if !contains(validEnvs, strings.ToLower(cfg.Environment)) {
		errs = append(errs, fmt.Errorf("invalid environment: %s (must be one of: %s)",
			cfg.Environment, strings.Join(validEnvs, ", ")))
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

// validateServer validates server configuration
func validateServer(cfg *ServerConfig) error {
	var errs []error

	if cfg.HTTPPort == "" {
		errs = append(errs, errors.New("HTTP port is required"))
	}

	if cfg.GRPCPort == "" {
		errs = append(errs, errors.New("gRPC port is required"))
	}

	if cfg.ReadTimeout <= 0 {
		errs = append(errs, errors.New("read timeout must be positive"))
	}

	if cfg.WriteTimeout <= 0 {
		errs = append(errs, errors.New("write timeout must be positive"))
	}

	if cfg.IdleTimeout <= 0 {
		errs = append(errs, errors.New("idle timeout must be positive"))
	}

	if cfg.GracefulShutdownTimeout <= 0 {
		errs = append(errs, errors.New("graceful shutdown timeout must be positive"))
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

// validateDatabase validates database configuration
func validateDatabase(cfg *DatabaseConfig) error {
	var errs []error

	if cfg.URI == "" {
		errs = append(errs, errors.New("database URI is required"))
	}

	if cfg.Name == "" {
		errs = append(errs, errors.New("database name is required"))
	}

	if cfg.MaxPoolSize == 0 {
		errs = append(errs, errors.New("max pool size must be positive"))
	}

	if cfg.MinPoolSize > cfg.MaxPoolSize {
		errs = append(errs, errors.New("min pool size cannot exceed max pool size"))
	}

	if cfg.ConnectionTimeout <= 0 {
		errs = append(errs, errors.New("connection timeout must be positive"))
	}

	if cfg.QueryTimeout <= 0 {
		errs = append(errs, errors.New("query timeout must be positive"))
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

// validateJWT validates JWT configuration
func validateJWT(cfg *JWTConfig) error {
	var errs []error

	if cfg.SecretKey == "" {
		errs = append(errs, errors.New("JWT secret key is required"))
	}

	// Warn about insecure secret in production
	if cfg.SecretKey == "change-me-in-production" {
		errs = append(errs, errors.New("JWT secret key must be changed from default value"))
	}

	// Enforce minimum secret length
	if len(cfg.SecretKey) < 32 {
		errs = append(errs, fmt.Errorf("JWT secret key too short (got %d characters, need at least 32)",
			len(cfg.SecretKey)))
	}

	if cfg.AccessTokenExpiry <= 0 {
		errs = append(errs, errors.New("access token expiry must be positive"))
	}

	if cfg.RefreshTokenExpiry <= 0 {
		errs = append(errs, errors.New("refresh token expiry must be positive"))
	}

	// Ensure refresh token lives longer than access token
	if cfg.RefreshTokenExpiry <= cfg.AccessTokenExpiry {
		errs = append(errs, errors.New("refresh token expiry must be greater than access token expiry"))
	}

	if cfg.Issuer == "" {
		errs = append(errs, errors.New("JWT issuer is required"))
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

// validateLogger validates logger configuration
func validateLogger(cfg *LoggerConfig) error {
	var errs []error

	validLevels := []string{"debug", "info", "warn", "error"}
	if !contains(validLevels, strings.ToLower(cfg.Level)) {
		errs = append(errs, fmt.Errorf("invalid log level: %s (must be one of: %s)",
			cfg.Level, strings.Join(validLevels, ", ")))
	}

	validFormats := []string{"json", "text"}
	if !contains(validFormats, strings.ToLower(cfg.Format)) {
		errs = append(errs, fmt.Errorf("invalid log format: %s (must be one of: %s)",
			cfg.Format, strings.Join(validFormats, ", ")))
	}

	if cfg.Output == "" {
		errs = append(errs, errors.New("log output is required"))
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

// contains checks if slice contains string
func contains(slice []string, str string) bool {
	for _, s := range slice {
		if strings.EqualFold(s, str) {
			return true
		}
	}
	return false
}

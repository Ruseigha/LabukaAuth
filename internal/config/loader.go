package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func Load() (*Config, error) {
	// Start with defaults
	cfg := DefaultConfig()

	// Override with environment variables
	// App config
	if v := os.Getenv("APP_NAME"); v != "" {
		cfg.App.Name = v
	}
	if v := os.Getenv("APP_ENV"); v != "" {
		cfg.App.Environment = v
	}
	if v := os.Getenv("APP_VERSION"); v != "" {
		cfg.App.Version = v
	}
	if v := os.Getenv("APP_DEBUG"); v != "" {
		cfg.App.Debug = parseBool(v)
	}

	// Server config
	if v := os.Getenv("HTTP_PORT"); v != "" {
		cfg.Server.HTTPPort = v
	}
	if v := os.Getenv("GRPC_PORT"); v != "" {
		cfg.Server.GRPCPort = v
	}
	if v := os.Getenv("SERVER_READ_TIMEOUT"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			cfg.Server.ReadTimeout = d
		}
	}
	if v := os.Getenv("SERVER_WRITE_TIMEOUT"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			cfg.Server.WriteTimeout = d
		}
	}
	if v := os.Getenv("SERVER_IDLE_TIMEOUT"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			cfg.Server.IdleTimeout = d
		}
	}
	if v := os.Getenv("SERVER_GRACEFUL_SHUTDOWN_TIMEOUT"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			cfg.Server.GracefulShutdownTimeout = d
		}
	}

	// Database config
	if v := os.Getenv("MONGO_URI"); v != "" {
		cfg.Database.URI = v
	}
	if v := os.Getenv("MONGO_DATABASE"); v != "" {
		cfg.Database.Name = v
	}
	if v := os.Getenv("MONGO_MAX_POOL_SIZE"); v != "" {
		if n, err := strconv.ParseUint(v, 10, 64); err == nil {
			cfg.Database.MaxPoolSize = n
		}
	}
	if v := os.Getenv("MONGO_MIN_POOL_SIZE"); v != "" {
		if n, err := strconv.ParseUint(v, 10, 64); err == nil {
			cfg.Database.MinPoolSize = n
		}
	}
	if v := os.Getenv("MONGO_CONNECTION_TIMEOUT"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			cfg.Database.ConnectionTimeout = d
		}
	}
	if v := os.Getenv("MONGO_QUERY_TIMEOUT"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			cfg.Database.QueryTimeout = d
		}
	}

	// JWT config
	if v := os.Getenv("JWT_SECRET_KEY"); v != "" {
		cfg.JWT.SecretKey = v
	}
	if v := os.Getenv("JWT_ACCESS_TOKEN_EXPIRY"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			cfg.JWT.AccessTokenExpiry = d
		}
	}
	if v := os.Getenv("JWT_REFRESH_TOKEN_EXPIRY"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			cfg.JWT.RefreshTokenExpiry = d
		}
	}
	if v := os.Getenv("JWT_ISSUER"); v != "" {
		cfg.JWT.Issuer = v
	}

	// Logger config
	if v := os.Getenv("LOG_LEVEL"); v != "" {
		cfg.Logger.Level = v
	}
	if v := os.Getenv("LOG_FORMAT"); v != "" {
		cfg.Logger.Format = v
	}
	if v := os.Getenv("LOG_OUTPUT"); v != "" {
		cfg.Logger.Output = v
	}

	// Validate configuration
	if err := Validate(cfg); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

func LoadFromFile(filepath string) (*Config, error) {
	// Load .env file
	if err := loadEnvFile(filepath); err != nil {
		// .env file is optional - don't fail if missing
		// WHY: In production, use real env vars, not .env file
	}

	// Load from environment (potentially overriding .env values)
	return Load()
}

func loadEnvFile(filepath string) error {
	// Read file
	data, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	// Parse lines
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		// Skip empty lines and comments
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse KEY=VALUE
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Remove quotes if present
		value = strings.Trim(value, `"'`)

		// Set environment variable (only if not already set)
		// WHY: Real env vars take precedence over .env file
		if os.Getenv(key) == "" {
			os.Setenv(key, value)
		}
	}

	return nil
}

// Accepts: true, false, 1, 0, yes, no (case-insensitive)
func parseBool(s string) bool {
	s = strings.ToLower(strings.TrimSpace(s))
	return s == "true" || s == "1" || s == "yes"
}

// MustLoad loads configuration and panics on error
func MustLoad() *Config {
	cfg, err := Load()
	if err != nil {
		panic(fmt.Sprintf("Failed to load configuration: %v", err))
	}
	return cfg
}

// MustLoadFromFile loads from file and panics on validation error
func MustLoadFromFile(filepath string) *Config {
	cfg, err := LoadFromFile(filepath)
	if err != nil {
		panic(fmt.Sprintf("Failed to load configuration from %s: %v", filepath, err))
	}
	return cfg
}

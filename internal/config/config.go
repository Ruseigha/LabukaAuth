package config

import "time"

type Config struct {
	// App contains application-level settings
	App AppConfig

	// Server contains HTTP/gRPC server settings
	Server ServerConfig

	// Database contains MongoDB configuration
	Database DatabaseConfig

	// JWT contains authentication settings
	JWT JWTConfig

	// Logger contains logging configuration
	Logger LoggerConfig
}

type AppConfig struct {
	Name        string
	Environment string
	Version     string
	Debug       bool
}

type ServerConfig struct {
	HTTPPort                string
	GRPCPort                string
	ReadTimeout             time.Duration
	WriteTimeout            time.Duration
	IdleTimeout             time.Duration
	GracefulShutdownTimeout time.Duration
}

type DatabaseConfig struct {
	URI               string
	Name              string
	MaxPoolSize       uint64
	MinPoolSize       uint64
	ConnectionTimeout time.Duration
	QueryTimeout      time.Duration
}

type JWTConfig struct {
	SecretKey          string
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
	Issuer             string
}

type LoggerConfig struct {
	Level  string
	Format string
	Output string
}

func DefaultConfig() *Config {
	return &Config{
		App: AppConfig{
			Name:        "auth-service",
			Environment: "development",
			Version:     "1.0.0",
			Debug:       true,
		},
		Server: ServerConfig{
			HTTPPort:                "8001",
			GRPCPort:                "9001",
			ReadTimeout:             10 * time.Second,
			WriteTimeout:            10 * time.Second,
			IdleTimeout:             60 * time.Second,
			GracefulShutdownTimeout: 30 * time.Second,
		},
		Database: DatabaseConfig{
			URI:               "mongodb://localhost:27017",
			Name:              "auth_service_dev",
			MaxPoolSize:       100,
			MinPoolSize:       10,
			ConnectionTimeout: 10 * time.Second,
			QueryTimeout:      5 * time.Second,
		},
		JWT: JWTConfig{
			SecretKey:          "change-me-in-production", // WARNING: Not secure!
			AccessTokenExpiry:  15 * time.Minute,
			RefreshTokenExpiry: 7 * 24 * time.Hour, // 7 days
			Issuer:             "auth-service",
		},
		Logger: LoggerConfig{
			Level:  "debug",
			Format: "text",
			Output: "stdout",
		},
	}
}

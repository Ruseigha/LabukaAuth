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
	HTTPPort string
	GRPCPort string
	ReadTimeout time.Duration
	WriteTimeout time.Duration
	IdleTimeout time.Duration
	GracefulShutdownTimeout time.Duration
}


type DatabaseConfig struct {
  URI string
  Name string
  MaxPoolSize uint64
  MinPoolSize uint64
  ConnectionTimeout time.Duration
  QueryTimeout time.Duration
}


type JWTConfig struct {
  SecretKey string
  AccessTokenExpiry time.Duration
  RefreshTokenExpiry time.Duration
  Issuer string
}

type LoggerConfig struct {
  Level string
  Format string
  Output string
}
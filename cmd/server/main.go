package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Ruseigha/LabukaAuth/internal/config"
	httpdelivery "github.com/Ruseigha/LabukaAuth/internal/delivery/http"
	"github.com/Ruseigha/LabukaAuth/internal/infrastructure/persistence/mongodb"
	"github.com/Ruseigha/LabukaAuth/internal/infrastructure/security"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Printf("Starting %s v%s in %s mode", cfg.App.Name, cfg.App.Version, cfg.App.Environment)

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	mongoClient, err := mongodb.NewClient(ctx, cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoClient.Close(context.Background())

	log.Println("✓ Connected to MongoDB")

	// Create indexes
	if err := mongodb.CreateIndexes(ctx, mongoClient.Collection("users")); err != nil {
		log.Fatalf("Failed to create indexes: %v", err)
	}

	log.Println("✓ Database indexes created")

	// Initialize infrastructure
	userRepo := mongodb.NewUserRepository(mongoClient.Database())
	passwordHasher := security.NewBcryptHasher(10) // Cost factor 10
	jwtGenerator := security.NewJWTGenerator(
		cfg.JWT.SecretKey,
		cfg.JWT.AccessTokenExpiry,
		cfg.JWT.RefreshTokenExpiry,
		cfg.JWT.Issuer,
	)

	// Initialize use cases
	authService := auth.NewAuthService(userRepo, passwordHasher, jwtGenerator)

	// Setup HTTP router
	router := httpdelivery.SetupRouter(authService, cfg.App.Version)

	// Create HTTP server
	server := &http.Server{
		Addr:         ":" + cfg.Server.HTTPPort,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Start server in goroutine
	go func() {
		log.Printf("🚀 HTTP server listening on port %s", cfg.Server.HTTPPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server...")

	// Graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), cfg.Server.GracefulShutdownTimeout)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("✓ Server stopped gracefully")
}

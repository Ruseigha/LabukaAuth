package testutil

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/Ruseigha/LabukaAuth/internal/config"
	"github.com/Ruseigha/LabukaAuth/internal/infrastructure/persistence/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
)

type TestMongoDBClient struct {
	client     *mongodb.Client
	database   *mongo.Database
	collection *mongo.Collection
}

func SetupMongoDB(t *testing.T) *TestMongoDBClient {
	t.Helper()

	// Load test configuration
	cfg := config.DefaultConfig()
	cfg.Database.Name = "auth_service_test"
	cfg.Database.URI = "mongodb://localhost:27017"

	// Override from environment if set
	// WHY: CI/CD might use different host
	if uri := os.Getenv("MONGO_URI"); uri != "" {
		cfg.Database.URI = uri
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongodb.NewClient(ctx, cfg.Database)
	if err != nil {
		t.Fatalf("Failed to connect to test MongoDB: %v", err)
	}

	// Get database and collection
	db := client.Database()
	collection := db.Collection("users")

	// Clean database before tests
	// WHY: Start with clean slate
	if err := cleanDatabase(ctx, db); err != nil {
		t.Fatalf("Failed to clean test database: %v", err)
	}

	// Create indexes
	if err := mongodb.CreateIndexes(ctx, collection); err != nil {
		t.Fatalf("Failed to create indexes: %v", err)
	}

	return &TestMongoDBClient{
		client:     client,
		database:   db,
		collection: collection,
	}
}

func (tc *TestMongoDBClient) Cleanup(t *testing.T) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Clean database
	if err := cleanDatabase(ctx, tc.database); err != nil {
		t.Logf("Warning: Failed to clean database: %v", err)
	}

	// Close connection
	if err := tc.client.Close(ctx); err != nil {
		t.Logf("Warning: Failed to close MongoDB connection: %v", err)
	}
}

// Client returns the MongoDB client
func (tc *TestMongoDBClient) Client() *mongodb.Client {
	return tc.client
}

// Database returns the MongoDB database
func (tc *TestMongoDBClient) Database() *mongo.Database {
	return tc.database
}

// Collection returns the users collection
func (tc *TestMongoDBClient) Collection() *mongo.Collection {
	return tc.collection
}

// CleanCollection removes all documents from users collection
// WHY: Clean state between test cases
func (tc *TestMongoDBClient) CleanCollection(t *testing.T) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := tc.collection.DeleteMany(ctx, map[string]interface{}{})
	if err != nil {
		t.Fatalf("Failed to clean collection: %v", err)
	}
}

// cleanDatabase drops all collections in database
func cleanDatabase(ctx context.Context, db *mongo.Database) error {
	// Get all collection names
	collections, err := db.ListCollectionNames(ctx, map[string]interface{}{})
	if err != nil {
		return fmt.Errorf("failed to list collections: %w", err)
	}

	// Drop each collection
	for _, name := range collections {
		if err := db.Collection(name).Drop(ctx); err != nil {
			return fmt.Errorf("failed to drop collection %s: %w", name, err)
		}
	}

	return nil
}

// WaitForMongoDB waits until MongoDB is ready
// WHY: Useful for CI/CD where MongoDB might still be starting
func WaitForMongoDB(uri string, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cfg := config.DatabaseConfig{
		URI:               uri,
		Name:              "test",
		MaxPoolSize:       10,
		MinPoolSize:       2,
		ConnectionTimeout: 5 * time.Second,
		QueryTimeout:      5 * time.Second,
	}

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for MongoDB")
		case <-ticker.C:
			client, err := mongodb.NewClient(ctx, cfg)
			if err == nil {
				client.Close(ctx)
				return nil
			}
		}
	}
}

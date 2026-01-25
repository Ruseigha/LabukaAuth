package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/Ruseigha/LabukaAuth/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Client struct {
	client   *mongo.Client
	database *mongo.Database
}

func NewClient(ctx context.Context, cfg config.DatabaseConfig) (*Client, error) {
	// Validate config
	if cfg.URI == "" {
		return nil, fmt.Errorf("MongoDB URI is required")
	}

	if cfg.Name == "" {
		return nil, fmt.Errorf("MongoDB database name is required")
	}

	// Create client options
	clientOpts := options.Client().
		ApplyURI(cfg.URI).
		SetMaxPoolSize(cfg.MaxPoolSize).
		SetMinPoolSize(cfg.MinPoolSize).
		SetConnectTimeout(cfg.ConnectionTimeout).
		SetServerSelectionTimeout(cfg.ConnectionTimeout)

	// Connect
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping to verify connection
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := client.Ping(pingCtx, readpref.Primary()); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	return &Client{
		client:   client,
		database: client.Database(cfg.Name),
	}, nil
}

// Database returns the MongoDB database
func (c *Client) Database() *mongo.Database {
	return c.database
}

// Collection returns a collection
func (c *Client) Collection(name string) *mongo.Collection {
	return c.database.Collection(name)
}

// Ping checks MongoDB connectivity
func (c *Client) Ping(ctx context.Context) error {
	return c.client.Ping(ctx, readpref.Primary())
}

// Close disconnects from MongoDB
func (c *Client) Close(ctx context.Context) error {
	if c.client == nil {
		return nil
	}
	return c.client.Disconnect(ctx)
}

// IsConnected checks if connected
func (c *Client) IsConnected(ctx context.Context) bool {
	if c.client == nil {
		return false
	}
	return c.Ping(ctx) == nil
}

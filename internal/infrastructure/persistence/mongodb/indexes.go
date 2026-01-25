package mongodb

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateIndexes(ctx context.Context, collection *mongo.Collection) error {
	emailIndexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "email", Value: 1}, // 1 = ascending order
		},
		Options: options.Index().
			SetUnique(true).             // Enforce uniqueness
			SetName("email_unique_idx"), // Name for management
	}

	// Created_at index - for sorting/filtering
	// WHY: Common query: Find users ordered by creation date
	createdAtIndexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "created_at", Value: -1}, // -1 = descending (newest first)
		},
		Options: options.Index().
			SetName("created_at_idx"),
	}

	// Is_active index - for filtering active users
	// WHY: Common query: Find all active users
	isActiveIndexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "is_active", Value: 1},
		},
		Options: options.Index().
			SetName("is_active_idx"),
	}

	compoundIndexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "is_active", Value: 1},
			{Key: "created_at", Value: -1},
		},
		Options: options.Index().
			SetName("is_active_created_at_idx"),
	}

	// Create all indexes
	indexModels := []mongo.IndexModel{
		emailIndexModel,
		createdAtIndexModel,
		isActiveIndexModel,
		compoundIndexModel,
	}

	_, err := collection.Indexes().CreateMany(ctx, indexModels)
	if err != nil {
		return fmt.Errorf("failed to create indexes: %w", err)
	}

	return nil
}

func DropIndexes(ctx context.Context, collection *mongo.Collection) error {
	// Get all index names
	cursor, err := collection.Indexes().List(ctx)
	if err != nil {
		return fmt.Errorf("failed to list indexes: %w", err)
	}
	defer cursor.Close(ctx)

	// Drop each index (except _id which is mandatory)
	for cursor.Next(ctx) {
		var index bson.M
		if err := cursor.Decode(&index); err != nil {
			continue
		}

		indexName, ok := index["name"].(string)
		if !ok || indexName == "_id_" {
			continue // Skip _id index (can't drop it)
		}

		if _, err := collection.Indexes().DropOne(ctx, indexName); err != nil {
			return fmt.Errorf("failed to drop index %s: %w", indexName, err)
		}
	}

	return nil
}
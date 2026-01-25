package mongodb

import (
	"context"
	"errors"
	"fmt"

	"github.com/Ruseigha/LabukaAuth/internal/domain/entity"
	"github.com/Ruseigha/LabukaAuth/internal/domain/valueobject"
	"github.com/Ruseigha/LabukaAuth/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	collection := db.Collection("users")
	return &UserRepository{
		collection: collection,
	}
}

func (r *UserRepository) Create(ctx context.Context, user *entity.User) error {
	// Convert domain entity to MongoDB document
	doc := fromEntity(user)

	// Insert into MongoDB
	_, err := r.collection.InsertOne(ctx, doc)
	if err != nil {
		// Check if error is duplicate key (unique constraint violation)
		if mongo.IsDuplicateKeyError(err) {
			// Email already exists (unique index violated)
			return repository.NewUserAlreadyExistsError("Create", err)
		}

		// Other database errors
		return repository.NewDatabaseQueryError("Create", err)
	}

	return nil
}

func (r *UserRepository) FindByID(ctx context.Context, id valueobject.UserID) (*entity.User, error) {
	// Build query filter
	filter := bson.M{"_id": id.String()}

	// Execute query
	var doc UserDocument
	err := r.collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		// Check if user not found
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, repository.NewUserNotFoundError("FindByID")
		}

		// Other errors (connection, timeout, etc.)
		return nil, repository.NewDatabaseQueryError("FindByID", err)
	}

	// Convert MongoDB document to domain entity
	user, err := doc.toEntity()
	if err != nil {
		// Invalid data in database (corrupted)
		return nil, repository.NewDatabaseQueryError("FindByID", fmt.Errorf("invalid user data: %w", err))
	}

	return user, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email valueobject.Email) (*entity.User, error) {

	filter := bson.M{"email": email.String()}

	// Execute query
	var doc UserDocument
	err := r.collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, repository.NewUserNotFoundError("FindByEmail")
		}
		return nil, repository.NewDatabaseQueryError("FindByEmail", err)
	}

	// Convert to entity
	user, err := doc.toEntity()
	if err != nil {
		return nil, repository.NewDatabaseQueryError("FindByEmail", fmt.Errorf("invalid user data: %w", err))
	}

	return user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *entity.User) error {
	// Build filter (find by ID)
	filter := bson.M{"_id": user.ID().String()}

	// Build update document
	update := bson.M{
		"$set": bson.M{
			"email":      user.Email().String(),
			"password":   user.Password().Hash(),
			"updated_at": user.UpdatedAt(),
			"is_active":  user.IsActive(),
			// Note: Don't update created_at (immutable)
		},
	}

	// Execute update
	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		// Check for duplicate email (if email was changed)
		if mongo.IsDuplicateKeyError(err) {
			return repository.NewUserAlreadyExistsError("Update", err)
		}
		return repository.NewDatabaseQueryError("Update", err)
	}

	// Check if user was found
	if result.MatchedCount == 0 {
		return repository.NewUserNotFoundError("Update")
	}

	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id valueobject.UserID) error {
	// Build filter
	filter := bson.M{"_id": id.String()}

	// Execute delete
	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return repository.NewDatabaseQueryError("Delete", err)
	}

	// Check if user was found
	if result.DeletedCount == 0 {
		return repository.NewUserNotFoundError("Delete")
	}

	return nil
}

func (r *UserRepository) ExistsByEmail(ctx context.Context, email valueobject.Email) (bool, error) {
	// Build filter
	filter := bson.M{"email": email.String()}

	// Count documents (more efficient than FindOne)
	// WHY: We only need existence, not the actual document
	count, err := r.collection.CountDocuments(ctx, filter, options.Count().SetLimit(1))
	if err != nil {
		return false, repository.NewDatabaseQueryError("ExistsByEmail", err)
	}

	return count > 0, nil
}

func (r *UserRepository) List(ctx context.Context, offset, limit int) ([]*entity.User, error) {
	// Validate pagination parameters
	if offset < 0 {
		offset = 0
	}
	if limit <= 0 || limit > 100 {
		limit = 20 // Default limit
	}

	// Build query options
	findOptions := options.Find().
		SetSkip(int64(offset)).
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	// Execute query
	cursor, err := r.collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, repository.NewDatabaseQueryError("List", err)
	}
	defer cursor.Close(ctx)

	// Decode results
	var users []*entity.User
	for cursor.Next(ctx) {
		var doc UserDocument
		if err := cursor.Decode(&doc); err != nil {
			// Log error but continue processing other documents
			continue
		}

		user, err := doc.toEntity()
		if err != nil {
			// Log error but continue
			continue
		}

		users = append(users, user)
	}

	// Check for cursor errors
	if err := cursor.Err(); err != nil {
		return nil, repository.NewDatabaseQueryError("List", err)
	}

	return users, nil
}

func (r *UserRepository) Count(ctx context.Context) (int64, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, repository.NewDatabaseQueryError("Count", err)
	}
	return count, nil
}

package mongodb_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Ruseigha/LabukaAuth/internal/domain/valueobject"
	mongodbpkg "github.com/Ruseigha/LabukaAuth/internal/infrastructure/persistence/mongodb"
	"github.com/Ruseigha/LabukaAuth/internal/repository"
	"github.com/Ruseigha/LabukaAuth/test/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserRepository_Create(t *testing.T) {
	// Setup
	testDB := testutil.SetupMongoDB(t)
	defer testDB.Cleanup(t)

	repo := mongodbpkg.NewUserRepository(testDB.Database())
	ctx := context.Background()

	t.Run("success - create new user", func(t *testing.T) {
		// Arrange
		user := testutil.CreateTestUser(t, "user1@example.com", "SecureP@ss123")

		// Act
		err := repo.Create(ctx, user)

		// Assert
		require.NoError(t, err)

		// Verify user exists in database
		found, err := repo.FindByID(ctx, user.ID())
		require.NoError(t, err)
		assert.Equal(t, user.Email().String(), found.Email().String())
	})

	t.Run("error - duplicate email", func(t *testing.T) {
		// Arrange
		testDB.CleanCollection(t)
		user1 := testutil.CreateTestUser(t, "duplicate@example.com", "SecureP@ss123")
		user2 := testutil.CreateTestUser(t, "duplicate@example.com", "DifferentP@ss456")

		// Act
		err1 := repo.Create(ctx, user1)
		err2 := repo.Create(ctx, user2)

		// Assert
		require.NoError(t, err1, "First user should be created")
		require.Error(t, err2, "Second user with same email should fail")
		assert.True(t, errors.Is(err2, repository.ErrUserAlreadyExists))
	})
}

// TestUserRepository_FindByID tests finding users by ID
func TestUserRepository_FindByID(t *testing.T) {
	testDB := testutil.SetupMongoDB(t)
	defer testDB.Cleanup(t)

	repo := mongodbpkg.NewUserRepository(testDB.Database())
	ctx := context.Background()

	t.Run("success - user found", func(t *testing.T) {
		// Arrange
		testDB.CleanCollection(t)
		user := testutil.CreateTestUser(t, "findme@example.com", "SecureP@ss123")
		err := repo.Create(ctx, user)
		require.NoError(t, err)

		// Act
		found, err := repo.FindByID(ctx, user.ID())

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, user.ID().String(), found.ID().String())
		assert.Equal(t, user.Email().String(), found.Email().String())
		assert.True(t, found.IsActive())
	})

	t.Run("error - user not found", func(t *testing.T) {
		// Arrange
		testDB.CleanCollection(t)
		nonExistentID := valueobject.NewUserID()

		// Act
		found, err := repo.FindByID(ctx, nonExistentID)

		// Assert
		require.Error(t, err)
		assert.Nil(t, found)
		assert.True(t, errors.Is(err, repository.ErrUserNotFound))
	})
}

// TestUserRepository_FindByEmail tests finding users by email
func TestUserRepository_FindByEmail(t *testing.T) {
	testDB := testutil.SetupMongoDB(t)
	defer testDB.Cleanup(t)

	repo := mongodbpkg.NewUserRepository(testDB.Database())
	ctx := context.Background()

	t.Run("success - user found", func(t *testing.T) {
		// Arrange
		testDB.CleanCollection(t)
		email := "findbyme@example.com"
		user := testutil.CreateTestUser(t, email, "SecureP@ss123")
		err := repo.Create(ctx, user)
		require.NoError(t, err)

		emailVO, _ := valueobject.NewEmail(email)

		// Act
		found, err := repo.FindByEmail(ctx, emailVO)

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, user.ID().String(), found.ID().String())
		assert.Equal(t, email, found.Email().String())
	})

	t.Run("success - case insensitive search", func(t *testing.T) {
		// Arrange
		testDB.CleanCollection(t)
		user := testutil.CreateTestUser(t, "CaseTest@Example.COM", "SecureP@ss123")
		err := repo.Create(ctx, user)
		require.NoError(t, err)

		// Search with different case
		emailVO, _ := valueobject.NewEmail("casetest@example.com")

		// Act
		found, err := repo.FindByEmail(ctx, emailVO)

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, found)
	})

	t.Run("error - user not found", func(t *testing.T) {
		// Arrange
		testDB.CleanCollection(t)
		emailVO, _ := valueobject.NewEmail("notfound@example.com")

		// Act
		found, err := repo.FindByEmail(ctx, emailVO)

		// Assert
		require.Error(t, err)
		assert.Nil(t, found)
		assert.True(t, errors.Is(err, repository.ErrUserNotFound))
	})
}

// TestUserRepository_Update tests updating users
func TestUserRepository_Update(t *testing.T) {
	testDB := testutil.SetupMongoDB(t)
	defer testDB.Cleanup(t)

	repo := mongodbpkg.NewUserRepository(testDB.Database())
	ctx := context.Background()

	t.Run("success - update email", func(t *testing.T) {
		// Arrange
		testDB.CleanCollection(t)
		user := testutil.CreateTestUser(t, "original@example.com", "SecureP@ss123")
		err := repo.Create(ctx, user)
		require.NoError(t, err)

		// Update email
		newEmail, _ := valueobject.NewEmail("updated@example.com")
		err = user.UpdateEmail(newEmail)
		require.NoError(t, err)

		// Act
		err = repo.Update(ctx, user)

		// Assert
		require.NoError(t, err)

		// Verify update
		found, err := repo.FindByID(ctx, user.ID())
		require.NoError(t, err)
		assert.Equal(t, "updated@example.com", found.Email().String())
	})

	t.Run("success - deactivate user", func(t *testing.T) {
		// Arrange
		testDB.CleanCollection(t)
		user := testutil.CreateTestUser(t, "deactivate@example.com", "SecureP@ss123")
		err := repo.Create(ctx, user)
		require.NoError(t, err)

		// Deactivate
		user.Deactivate()

		// Act
		err = repo.Update(ctx, user)

		// Assert
		require.NoError(t, err)

		// Verify
		found, err := repo.FindByID(ctx, user.ID())
		require.NoError(t, err)
		assert.False(t, found.IsActive())
	})

	t.Run("error - user not found", func(t *testing.T) {
		// Arrange
		testDB.CleanCollection(t)
		user := testutil.CreateTestUser(t, "phantom@example.com", "SecureP@ss123")
		// Don't create user in database

		// Act
		err := repo.Update(ctx, user)

		// Assert
		require.Error(t, err)
		assert.True(t, errors.Is(err, repository.ErrUserNotFound))
	})
}

// TestUserRepository_Delete tests deleting users
func TestUserRepository_Delete(t *testing.T) {
	testDB := testutil.SetupMongoDB(t)
	defer testDB.Cleanup(t)

	repo := mongodbpkg.NewUserRepository(testDB.Database())
	ctx := context.Background()

	t.Run("success - delete existing user", func(t *testing.T) {
		// Arrange
		testDB.CleanCollection(t)
		user := testutil.CreateTestUser(t, "deleteme@example.com", "SecureP@ss123")
		err := repo.Create(ctx, user)
		require.NoError(t, err)

		// Act
		err = repo.Delete(ctx, user.ID())

		// Assert
		require.NoError(t, err)

		// Verify deleted
		found, err := repo.FindByID(ctx, user.ID())
		assert.Error(t, err)
		assert.Nil(t, found)
		assert.True(t, errors.Is(err, repository.ErrUserNotFound))
	})

	t.Run("error - user not found", func(t *testing.T) {
		// Arrange
		testDB.CleanCollection(t)
		nonExistentID := valueobject.NewUserID()

		// Act
		err := repo.Delete(ctx, nonExistentID)

		// Assert
		require.Error(t, err)
		assert.True(t, errors.Is(err, repository.ErrUserNotFound))
	})
}

// TestUserRepository_ExistsByEmail tests checking email existence
func TestUserRepository_ExistsByEmail(t *testing.T) {
	testDB := testutil.SetupMongoDB(t)
	defer testDB.Cleanup(t)

	repo := mongodbpkg.NewUserRepository(testDB.Database())
	ctx := context.Background()

	t.Run("returns true when user exists", func(t *testing.T) {
		// Arrange
		testDB.CleanCollection(t)
		email := "exists@example.com"
		user := testutil.CreateTestUser(t, email, "SecureP@ss123")
		err := repo.Create(ctx, user)
		require.NoError(t, err)

		emailVO, _ := valueobject.NewEmail(email)

		// Act
		exists, err := repo.ExistsByEmail(ctx, emailVO)

		// Assert
		require.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("returns false when user doesn't exist", func(t *testing.T) {
		// Arrange
		testDB.CleanCollection(t)
		emailVO, _ := valueobject.NewEmail("notexists@example.com")

		// Act
		exists, err := repo.ExistsByEmail(ctx, emailVO)

		// Assert
		require.NoError(t, err)
		assert.False(t, exists)
	})
}

// TestUserRepository_List tests listing users with pagination
func TestUserRepository_List(t *testing.T) {
	testDB := testutil.SetupMongoDB(t)
	defer testDB.Cleanup(t)

	repo := mongodbpkg.NewUserRepository(testDB.Database())
	ctx := context.Background()

	t.Run("success - list with pagination", func(t *testing.T) {
		// Arrange
		testDB.CleanCollection(t)

		// Create 5 users
		for i := 1; i <= 5; i++ {
			user := testutil.CreateTestUser(t, testutil.UniqueEmail("user"), "SecureP@ss123")
			err := repo.Create(ctx, user)
			require.NoError(t, err)
			time.Sleep(10 * time.Millisecond) // Ensure different timestamps
		}

		// Act - Get first 3
		users, err := repo.List(ctx, 0, 3)

		// Assert
		require.NoError(t, err)
		assert.Len(t, users, 3)

		// Act - Get next 2
		users, err = repo.List(ctx, 3, 3)

		// Assert
		require.NoError(t, err)
		assert.Len(t, users, 2)
	})

	t.Run("empty list when no users", func(t *testing.T) {
		// Arrange
		testDB.CleanCollection(t)

		// Act
		users, err := repo.List(ctx, 0, 10)

		// Assert
		require.NoError(t, err)
		assert.Empty(t, users)
	})
}

// TestUserRepository_Count tests counting users
func TestUserRepository_Count(t *testing.T) {
	testDB := testutil.SetupMongoDB(t)
	defer testDB.Cleanup(t)

	repo := mongodbpkg.NewUserRepository(testDB.Database())
	ctx := context.Background()

	t.Run("returns correct count", func(t *testing.T) {
		// Arrange
		testDB.CleanCollection(t)

		// Create 3 users
		for i := 1; i <= 3; i++ {
			user := testutil.CreateTestUser(t, testutil.UniqueEmail("count"), "SecureP@ss123")
			err := repo.Create(ctx, user)
			require.NoError(t, err)
		}

		// Act
		count, err := repo.Count(ctx)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, int64(3), count)
	})

	t.Run("returns zero when no users", func(t *testing.T) {
		// Arrange
		testDB.CleanCollection(t)

		// Act
		count, err := repo.Count(ctx)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, int64(0), count)
	})
}

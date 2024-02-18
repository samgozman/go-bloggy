package models

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

func NewTestDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	err = db.AutoMigrate(&User{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate: %w", err)
	}

	return db, nil
}

func testCreateUser(ctx context.Context, db *gorm.DB) (User, error) {
	user := User{
		ExternalID: uuid.New().String(),
		Login:      uuid.New().String(),
		AuthMethod: GitHubAuthMethod,
	}

	if err := db.WithContext(ctx).Create(&user).Error; err != nil {
		return User{}, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func TestUserDB(t *testing.T) {
	db, e := NewTestDB("users_test.db")
	assert.NoError(t, e)

	// Delete table after the test
	defer func() {
		_ = db.Migrator().DropTable(&User{})
	}()

	userDB := NewUserDB(db)

	t.Run("CreateUser", func(t *testing.T) {
		t.Run("should create user", func(t *testing.T) {
			user := &User{
				ExternalID: uuid.New().String(),
				Login:      uuid.New().String(),
				AuthMethod: GitHubAuthMethod,
			}

			err := userDB.CreateUser(context.Background(), user)
			assert.NoError(t, err)

			// Check that the user was created
			var u User
			err = userDB.conn.First(&u, user.ID).Error
			assert.NoError(t, err)
			assert.Equal(t, user.ID, u.ID)
			assert.Equal(t, user.ExternalID, u.ExternalID)
			assert.Equal(t, user.Login, u.Login)
			assert.Equal(t, user.AuthMethod, u.AuthMethod)
			assert.NotZero(t, u.CreatedAt)
			assert.NotZero(t, u.UpdatedAt)
		})

		t.Run("should return ErrDuplicate if user is duplicated", func(t *testing.T) {
			user := &User{
				ExternalID: uuid.New().String(),
				Login:      uuid.New().String(),
				AuthMethod: GitHubAuthMethod,
			}

			err := userDB.CreateUser(context.Background(), user)
			assert.NoError(t, err)

			err = userDB.CreateUser(context.Background(), user)
			assert.Error(t, err)
			fmt.Printf("Error: %v\n", err)
			assert.ErrorIs(t, err, ErrDuplicate)
		})
	})

	t.Run("GetUserByExternalID", func(t *testing.T) {
		t.Run("should get user by external id", func(t *testing.T) {
			ctx := context.Background()
			u, err := testCreateUser(ctx, userDB.conn)
			assert.NoError(t, err)

			user, err := userDB.GetUserByExternalID(ctx, u.ExternalID)
			assert.NoError(t, err)
			assert.Equal(t, u.ID, user.ID)
			assert.Equal(t, u.ExternalID, user.ExternalID)
			assert.Equal(t, u.Login, user.Login)
			assert.Equal(t, u.AuthMethod, user.AuthMethod)
		})

		t.Run("should return ErrNotFound if user is not found", func(t *testing.T) {
			user, err := userDB.GetUserByExternalID(context.Background(), uuid.New().String())
			assert.Error(t, err)
			assert.ErrorIs(t, err, ErrNotFound)
			assert.Nil(t, user)
		})
	})

	t.Run("GetUserByID", func(t *testing.T) {
		t.Run("should get user by id", func(t *testing.T) {
			ctx := context.Background()
			u, err := testCreateUser(ctx, userDB.conn)
			assert.NoError(t, err)

			user, err := userDB.GetUserByID(ctx, u.ID)
			assert.NoError(t, err)
			assert.Equal(t, u.ID, user.ID)
			assert.Equal(t, u.ExternalID, user.ExternalID)
			assert.Equal(t, u.Login, user.Login)
			assert.Equal(t, u.AuthMethod, user.AuthMethod)
		})

		t.Run("should return ErrNotFound if user is not found", func(t *testing.T) {
			user, err := userDB.GetUserByID(context.Background(), 0)
			assert.Error(t, err)
			assert.ErrorIs(t, err, ErrNotFound)
			assert.Nil(t, user)
		})
	})
}

func TestUser_Validate(t *testing.T) {
	t.Run("should return ErrUserLoginRequired if login is empty", func(t *testing.T) {
		user := User{
			ExternalID: uuid.New().String(),
			AuthMethod: GitHubAuthMethod,
		}
		err := user.Validate()
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrUserLoginRequired)
	})

	t.Run("should return ErrUserExternalIDRequired if external id is empty", func(t *testing.T) {
		user := User{
			Login:      uuid.New().String(),
			AuthMethod: GitHubAuthMethod,
		}
		err := user.Validate()
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrUserExternalIDRequired)
	})

	t.Run("should return ErrUserAuthMethodRequired if auth method is empty", func(t *testing.T) {
		user := User{
			Login:      uuid.New().String(),
			ExternalID: uuid.New().String(),
		}
		err := user.Validate()
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrUserAuthMethodRequired)
	})

	t.Run("should return nil if user is valid", func(t *testing.T) {
		user := User{
			Login:      uuid.New().String(),
			ExternalID: uuid.New().String(),
			AuthMethod: GitHubAuthMethod,
		}
		err := user.Validate()
		assert.NoError(t, err)
	})
}

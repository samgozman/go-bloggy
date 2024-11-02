package models

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	testdb "github.com/samgozman/go-bloggy/testutils/test-db"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"testing"
	"time"
)

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
	conn, err := testdb.InitDatabaseTest()
	assert.NoError(t, err)
	err = conn.AutoMigrate(&User{})
	assert.NoError(t, err)

	userDB := NewUserRepository(conn)

	t.Run("Upsert", func(t *testing.T) {
		t.Run("should create user", func(t *testing.T) {
			user := &User{
				ExternalID: uuid.New().String(),
				Login:      uuid.New().String(),
				AuthMethod: GitHubAuthMethod,
			}

			err := userDB.Upsert(context.Background(), user)
			assert.NoError(t, err)

			// Check that the user was created
			var u User
			err = userDB.conn.First(&u, user.ID).Error
			assert.NoError(t, err)
			assert.Equal(t, user.ID, u.ID)
			assert.Equal(t, user.ExternalID, u.ExternalID)
			assert.Equal(t, user.Login, u.Login)
			assert.Equal(t, user.AuthMethod, u.AuthMethod)
			assert.NotZero(t, u.CreatedAt.In(time.Local))
			assert.NotZero(t, u.UpdatedAt.In(time.Local))
		})

		t.Run("should update user", func(t *testing.T) {
			ctx := context.Background()
			u, err := testCreateUser(ctx, userDB.conn)
			assert.NoError(t, err)

			u.Login = uuid.New().String()
			err = userDB.Upsert(ctx, &u)
			assert.NoError(t, err)

			// Check that the user was updated
			var user User
			err = userDB.conn.First(&user, u.ID).Error
			assert.NoError(t, err)
			assert.Equal(t, u.ID, user.ID)
			assert.Equal(t, u.ExternalID, user.ExternalID)
			assert.Equal(t, u.Login, user.Login)
			assert.Equal(t, u.AuthMethod, user.AuthMethod)
			assert.Equal(t, u.CreatedAt.In(time.Local), user.CreatedAt.In(time.Local))
			assert.NotEqual(t, u.UpdatedAt.Truncate(time.Microsecond), user.UpdatedAt.Truncate(time.Microsecond))
		})

		t.Run("should return error if user is invalid", func(t *testing.T) {
			err := userDB.Upsert(context.Background(), &User{})
			assert.Error(t, err)
			assert.ErrorIs(t, err, ErrValidationFailed)
		})
	})

	t.Run("GetByExternalID", func(t *testing.T) {
		t.Run("should get user by external id", func(t *testing.T) {
			ctx := context.Background()
			u, err := testCreateUser(ctx, userDB.conn)
			assert.NoError(t, err)

			user, err := userDB.GetByExternalID(ctx, u.ExternalID)
			assert.NoError(t, err)
			assert.Equal(t, u.ID, user.ID)
			assert.Equal(t, u.ExternalID, user.ExternalID)
			assert.Equal(t, u.Login, user.Login)
			assert.Equal(t, u.AuthMethod, user.AuthMethod)
		})

		t.Run("should return ErrNotFound if user is not found", func(t *testing.T) {
			user, err := userDB.GetByExternalID(context.Background(), uuid.New().String())
			assert.Error(t, err)
			assert.ErrorIs(t, err, ErrNotFound)
			assert.Nil(t, user)
		})
	})

	t.Run("GetByID", func(t *testing.T) {
		t.Run("should get user by id", func(t *testing.T) {
			ctx := context.Background()
			u, err := testCreateUser(ctx, userDB.conn)
			assert.NoError(t, err)

			user, err := userDB.GetByID(ctx, u.ID)
			assert.NoError(t, err)
			assert.Equal(t, u.ID, user.ID)
			assert.Equal(t, u.ExternalID, user.ExternalID)
			assert.Equal(t, u.Login, user.Login)
			assert.Equal(t, u.AuthMethod, user.AuthMethod)
		})

		t.Run("should return ErrNotFound if user is not found", func(t *testing.T) {
			user, err := userDB.GetByID(context.Background(), 0)
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

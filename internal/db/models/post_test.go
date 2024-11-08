package models

import (
	"context"
	"github.com/google/uuid"
	testdb "github.com/samgozman/go-bloggy/testutils/test-db"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestPostDB(t *testing.T) {
	conn, err := testdb.InitDatabaseTest()
	assert.NoError(t, err)
	err = conn.AutoMigrate(&User{}, &Post{})
	assert.NoError(t, err)

	// insert a user to db
	user := &User{
		ExternalID: uuid.New().String(),
		Login:      uuid.New().String(),
		AuthMethod: GitHubAuthMethod,
	}
	err = conn.WithContext(context.Background()).Create(user).Error
	assert.NoError(t, err)

	postDB := NewPostRepository(conn)

	t.Run("Create", func(t *testing.T) {
		t.Run("create a new post", func(t *testing.T) {
			post := &Post{
				UserID:      user.ID,
				Slug:        uuid.New().String(),
				Title:       "Test Title",
				Description: "Test Description",
				Content:     "Test Content",
				Keywords:    "test,content",
			}

			err := postDB.Create(context.Background(), post)
			assert.NoError(t, err)
			assert.NotEmpty(t, post.ID)
			assert.NotZero(t, post.CreatedAt)
			assert.NotZero(t, post.UpdatedAt)
		})

		t.Run("return error if slug is not unique", func(t *testing.T) {
			post := &Post{
				UserID:      user.ID,
				Slug:        uuid.New().String(),
				Title:       "Test Title",
				Description: "Test Description",
				Content:     "Test Content",
			}

			err := postDB.Create(context.Background(), post)
			assert.NoError(t, err)

			err = postDB.Create(context.Background(), post)
			assert.Error(t, err)
			assert.ErrorIs(t, err, ErrDuplicate)
		})

		t.Run("return error if slug is not URL friendly", func(t *testing.T) {
			post := &Post{
				UserID:      user.ID,
				Slug:        "Test Title with spaces",
				Title:       "Test Title",
				Description: "Test Description",
				Content:     "Test Content",
			}

			err := postDB.Create(context.Background(), post)
			assert.Error(t, err)
			assert.ErrorIs(t, err, ErrValidationFailed)
			assert.ErrorIs(t, err, ErrPostInvalidSlug)
		})
	})

	t.Run("GetBySlug", func(t *testing.T) {
		anotherPost := &Post{
			UserID:      user.ID,
			Slug:        uuid.New().String(),
			Title:       "Test Title 00",
			Description: "Test Description",
			Content:     "Test Content",
		}

		post := &Post{
			UserID:      user.ID,
			Slug:        uuid.New().String(),
			Title:       "Test Title",
			Description: "Test Description",
			Content:     "Test Content",
		}

		t.Run("should get the post", func(t *testing.T) {
			err := postDB.Create(context.Background(), anotherPost)
			assert.NoError(t, err)

			err = postDB.Create(context.Background(), post)
			assert.NoError(t, err)

			retrievedPost, err := postDB.GetBySlug(context.Background(), post.Slug)
			assert.NoError(t, err)
			assert.Equal(t, post.Slug, retrievedPost.Slug)
		})

		t.Run("should return error if not found", func(t *testing.T) {
			_, err := postDB.GetBySlug(context.Background(), "not-found")
			assert.Error(t, err)
			assert.ErrorIs(t, err, ErrNotFound)
		})
	})

	t.Run("FindAll", func(t *testing.T) {
		for range 5 {
			post := &Post{
				UserID:              user.ID,
				Slug:                uuid.New().String(),
				Title:               "Test Title",
				Description:         "Test Description",
				Content:             "Test Content to read in 1 second",
				Keywords:            "some",
				SentToSubscribersAt: time.Now(),
			}

			err := postDB.Create(context.Background(), post)
			assert.NoError(t, err)
		}

		t.Run("get all", func(t *testing.T) {
			posts, err := postDB.FindAll(context.Background(), 1, 5)
			assert.NoError(t, err)
			assert.Equal(t, 5, len(posts))

			for i := 0; i < len(posts)-1; i++ { //nolint:intrange
				// check the order of the posts
				assert.True(t, posts[i].CreatedAt.After(posts[i+1].CreatedAt))
				assert.Equal(t, "Test Title", posts[i].Title)
				assert.Equal(t, "Test Description", posts[i].Description)
				assert.Equal(t, "some", posts[i].Keywords)
				assert.NotEmpty(t, posts[i].Slug)
				assert.NotEmpty(t, posts[i].CreatedAt)
				assert.Equal(t, 1, posts[i].ReadingTime)
				assert.NotZero(t, posts[i].SentToSubscribersAt)

				// check that unnecessary fields are empty
				assert.Empty(t, posts[i].Content)
				assert.Empty(t, posts[i].UpdatedAt)
				assert.Empty(t, posts[i].UserID)
			}
		})

		t.Run("get first 3", func(t *testing.T) {
			posts, err := postDB.FindAll(context.Background(), 1, 3)
			assert.NoError(t, err)
			assert.Equal(t, 3, len(posts))
		})

		t.Run("get next 3", func(t *testing.T) {
			posts, err := postDB.FindAll(context.Background(), 2, 3)
			assert.NoError(t, err)
			assert.Equal(t, 3, len(posts))
		})

		t.Run("return empty if non found", func(t *testing.T) {
			posts, err := postDB.FindAll(context.Background(), 1000, 1)
			assert.NoError(t, err)
			assert.Empty(t, posts)
		})
	})

	t.Run("Update", func(t *testing.T) {
		post := &Post{
			UserID:              user.ID,
			Slug:                uuid.New().String(),
			Title:               "Test Title",
			Description:         "Test Description",
			Content:             "Test Content",
			SentToSubscribersAt: time.Now(),
		}

		err := postDB.Create(context.Background(), post)
		assert.NoError(t, err)

		post.Title = "Updated Title"
		err = postDB.Update(context.Background(), post)
		assert.NoError(t, err)

		updatedPost, err := postDB.GetBySlug(context.Background(), post.Slug)
		assert.NoError(t, err)
		assert.Equal(t, "Updated Title", updatedPost.Title)
		assert.NotZero(t, updatedPost.UpdatedAt)
		assert.NotZero(t, updatedPost.SentToSubscribersAt)
	})
}

func TestPost_Validate(t *testing.T) {
	t.Run("return error if title is empty", func(t *testing.T) {
		post := &Post{
			Slug:        uuid.New().String(),
			Title:       "",
			Description: "Test Description",
			Content:     "Test Content",
			Keywords:    "test,content",
		}

		err := post.Validate()
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrPostTitleRequired)
	})

	t.Run("return error if description is empty", func(t *testing.T) {
		post := &Post{
			Slug:        uuid.New().String(),
			Title:       "Test",
			Description: "",
			Content:     "Test Content",
			Keywords:    "test,content",
		}

		err := post.Validate()
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrPostDescriptionRequired)
	})

	t.Run("return error if content is empty", func(t *testing.T) {
		post := &Post{
			Slug:        uuid.New().String(),
			Title:       "Test",
			Description: "Test Description",
			Content:     "",
			Keywords:    "test,content",
		}

		err := post.Validate()
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrPostContentRequired)
	})

	t.Run("return error if keywords is not in correct format", func(t *testing.T) {
		post := &Post{
			Slug:        uuid.New().String(),
			Title:       "Test",
			Description: "Test Description",
			Content:     "Test Content",
			Keywords:    "test,,content",
		}

		err := post.Validate()
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrPostWrongKeywordsString)
	})

	t.Run("return error if slug is not URL friendly", func(t *testing.T) {
		post := &Post{
			Slug:        "Test Title with spaces",
			Title:       "Test",
			Description: "Test Description",
			Content:     "Test Content",
			Keywords:    "test,content",
		}

		err := post.Validate()
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrPostInvalidSlug)
	})

	t.Run("return error if slug is empty", func(t *testing.T) {
		post := &Post{
			Slug:        "",
			Title:       "Test",
			Description: "Test Description",
			Content:     "Test Content",
			Keywords:    "test,content",
		}

		err := post.Validate()
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrPostURLRequired)
	})

	t.Run("return nil if valid", func(t *testing.T) {
		post := &Post{
			Slug:        uuid.New().String(),
			Title:       "Test",
			Description: "Test Description",
			Content:     "Test Content",
			Keywords:    "test,content",
		}

		err := post.Validate()
		assert.NoError(t, err)
	})
}

func TestPost_BeforeCreate(t *testing.T) {
	t.Run("return error if validation failed", func(t *testing.T) {
		post := &Post{
			Slug:        "",
			Title:       "Test",
			Description: "Test Description",
			Content:     "Test Content",
			Keywords:    "test,content",
		}

		err := post.BeforeCreate(nil)
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrValidationFailed)
	})

	t.Run("return nil if valid", func(t *testing.T) {
		post := &Post{
			Slug:        uuid.New().String(),
			Title:       "Test",
			Description: "Test Description",
			Content:     "Test Content text to read in about ~2 second",
			Keywords:    "test,content",
		}

		err := post.BeforeCreate(nil)
		assert.NoError(t, err)
		assert.Equal(t, 2, post.ReadingTime)
	})
}

func TestPost_BeforeUpdate(t *testing.T) {
	t.Run("return error if validation failed", func(t *testing.T) {
		post := &Post{
			Slug:        "",
			Title:       "Test",
			Description: "Test Description",
			Content:     "Test Content",
			Keywords:    "test,content",
		}

		err := post.BeforeUpdate(nil)
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrValidationFailed)
	})

	t.Run("return nil if valid", func(t *testing.T) {
		post := &Post{
			Slug:        uuid.New().String(),
			Title:       "Test",
			Description: "Test Description",
			Content:     "Test Content to read in 1 second",
			Keywords:    "test,content",
			ReadingTime: 5,
		}

		err := post.BeforeUpdate(nil)
		assert.NoError(t, err)
		assert.NotZero(t, post.UpdatedAt)
		assert.Equal(t, 1, post.ReadingTime) // reading time should be recalculated
	})
}

func TestPost_CountReadingTime(t *testing.T) {
	t.Run("return reading time", func(t *testing.T) {
		post := &Post{
			Content: "It should take 4 seconds to read with the average reading speed of 250 words per minute.",
		}

		assert.Equal(t, 4, int(post.CountReadingTime().Seconds()))
	})

	t.Run("return 0 if content is empty", func(t *testing.T) {
		post := &Post{
			Content: "",
		}

		assert.Zero(t, post.CountReadingTime().Seconds())
	})
}

func TestPostDB_Count(t *testing.T) {
	conn, err := testdb.InitDatabaseTest()
	assert.NoError(t, err)
	err = conn.AutoMigrate(&User{}, &Post{})
	assert.NoError(t, err)

	// insert a user to db
	user := &User{
		ExternalID: uuid.New().String(),
		Login:      uuid.New().String(),
		AuthMethod: GitHubAuthMethod,
	}

	err = conn.WithContext(context.Background()).Create(user).Error
	assert.NoError(t, err)

	postDB := NewPostRepository(conn)

	t.Run("return count of all posts", func(t *testing.T) {
		for range 5 {
			post := &Post{
				UserID:      user.ID,
				Slug:        uuid.New().String(),
				Title:       "Test Title",
				Description: "Test Description",
				Content:     "Test Content",
				Keywords:    "some",
			}

			err := postDB.Create(context.Background(), post)
			assert.NoError(t, err)
		}

		count, err := postDB.Count(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, int64(5), count)
	})
}

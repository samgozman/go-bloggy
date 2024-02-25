package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/testutil"
	"github.com/samgozman/go-bloggy/internal/db"
	"github.com/samgozman/go-bloggy/internal/db/models"
	"github.com/samgozman/go-bloggy/pkg/server"
	"github.com/stretchr/testify/assert"
	"net/http"
	"strconv"
	"testing"
)

func Test_PostPosts(t *testing.T) {
	conn, errDB := db.InitDatabase("file::memory:?cache=shared")
	if errDB != nil {
		t.Fatal(errDB)
	}

	// create user for test
	user := &models.User{
		ExternalID: uuid.New().String(),
		AuthMethod: models.GitHubAuthMethod,
		Login:      "testUser",
	}
	err := conn.Models.Users.Upsert(context.Background(), user)
	assert.NoError(t, err)

	jwtToken := "token"

	t.Run("OK", func(t *testing.T) {
		e, _, mockJwtService := registerHandlers(conn, []string{strconv.Itoa(user.ID)})
		mockJwtService.On("ParseTokenString", jwtToken).Return(user.ExternalID, nil)

		req := server.PostRequest{
			Title:       "Test Title",
			Slug:        uuid.New().String(),
			Content:     "Test Content",
			Description: "Test Description",
			Keywords:    &[]string{"test1", "test2"},
		}

		reqBody, _ := json.Marshal(req)

		res := testutil.NewRequest().
			Post("/posts").
			WithHeader("Content-Type", "application/json").
			WithBody(reqBody).
			WithJWSAuth(jwtToken).
			GoWithHTTPHandler(t, e)

		assert.Equal(t, http.StatusCreated, res.Code())

		var post server.PostResponse
		err := res.UnmarshalBodyToObject(&post)
		assert.NoError(t, err)

		assert.Equal(t, req.Title, post.Title)
		assert.Equal(t, req.Slug, post.Slug)
		assert.Equal(t, req.Content, post.Content)
		assert.Equal(t, req.Description, post.Description)
		assert.Equal(t, req.Keywords, post.Keywords)
		assert.NotEmpty(t, post.Id)
		assert.NotEmpty(t, post.CreatedAt)
		assert.NotEmpty(t, post.UpdatedAt)

		// check that post is in the database
		postFromDB, err := conn.Models.Posts.GetBySlug(context.Background(), req.Slug)
		assert.NoError(t, err)
		assert.Equal(t, req.Title, postFromDB.Title)
		assert.Equal(t, req.Slug, postFromDB.Slug)
		assert.Equal(t, req.Content, postFromDB.Content)
		assert.Equal(t, req.Description, postFromDB.Description)
		assert.Equal(t, postFromDB.Keywords, "test1,test2")
		assert.NotEmpty(t, postFromDB.CreatedAt)
		assert.NotEmpty(t, postFromDB.UpdatedAt)
	})

	t.Run("400 - errRequestBodyBinding - ErrUnsupportedMediaType", func(t *testing.T) {
		e, _, mockJwtService := registerHandlers(conn, nil)
		mockJwtService.On("ParseTokenString", jwtToken).Return(user.ExternalID, nil)

		req := server.PostRequest{
			Title:       "Test Title",
			Slug:        uuid.New().String(),
			Content:     "Test Content",
			Description: "Test Description",
			Keywords:    &[]string{"test1", "test2"},
		}

		reqBody, _ := json.Marshal(req)

		// Note: no Content-Type header
		res := testutil.NewRequest().
			Post("/posts").
			WithBody(reqBody).
			WithJWSAuth(jwtToken).
			GoWithHTTPHandler(t, e)

		assert.Equal(t, http.StatusBadRequest, res.Code())

		var body server.RequestError
		err := res.UnmarshalBodyToObject(&body)
		assert.NoError(t, err)
		assert.Equal(t, errRequestBodyBinding, body.Code)
		assert.Equal(t, body.Message, fmt.Sprintf("Error binding request body: %v", echo.ErrUnsupportedMediaType.Message))
	})

	t.Run("400 - errGetUser", func(t *testing.T) {
		e, _, mockJwtService := registerHandlers(conn, nil)
		mockJwtService.On("ParseTokenString", jwtToken).Return(uuid.New().String(), nil)

		req := server.PostRequest{
			Title:       "Test Title",
			Slug:        uuid.New().String(),
			Content:     "Test Content",
			Description: "Test Description",
		}

		reqBody, _ := json.Marshal(req)

		res := testutil.NewRequest().
			Post("/posts").
			WithHeader("Content-Type", "application/json").
			WithBody(reqBody).
			WithJWSAuth(jwtToken).
			GoWithHTTPHandler(t, e)

		assert.Equal(t, http.StatusBadRequest, res.Code())

		var body server.RequestError
		err := res.UnmarshalBodyToObject(&body)
		assert.NoError(t, err)
		assert.Equal(t, errGetUser, body.Code)
		assert.Equal(t, "Post author is not found", body.Message)
	})

	t.Run("409 - errDuplicatePost", func(t *testing.T) {
		e, _, mockJwtService := registerHandlers(conn, nil)
		mockJwtService.On("ParseTokenString", jwtToken).Return(user.ExternalID, nil)

		post1 := server.PostRequest{
			Title:       "Test Title",
			Slug:        uuid.New().String(),
			Content:     "Test Content",
			Description: "Test Description",
		}

		reqBody, _ := json.Marshal(post1)

		// create post with the same slug
		post2 := models.Post{
			UserID:      user.ID,
			Title:       post1.Title,
			Slug:        post1.Slug,
			Content:     post1.Content,
			Description: post1.Description,
		}

		err := conn.Models.Posts.Create(context.Background(), &post2)
		assert.NoError(t, err)

		res := testutil.NewRequest().
			Post("/posts").
			WithHeader("Content-Type", "application/json").
			WithBody(reqBody).
			WithJWSAuth(jwtToken).
			GoWithHTTPHandler(t, e)

		assert.Equal(t, http.StatusConflict, res.Code())

		var body server.RequestError
		err = res.UnmarshalBodyToObject(&body)
		assert.NoError(t, err)
		assert.Equal(t, errDuplicatePost, body.Code)
		assert.Equal(t, "Post with this URL slug already exists", body.Message)
	})

	t.Run("400 - errValidationFailed", func(t *testing.T) {
		e, _, mockJwtService := registerHandlers(conn, nil)
		mockJwtService.On("ParseTokenString", jwtToken).Return(user.ExternalID, nil)

		req := server.PostRequest{
			Title:       "Test Title",
			Slug:        uuid.New().String(),
			Content:     "Test Content",
			Description: "Test Description",
			Keywords:    &[]string{"test1", "", ""}, // invalid keywords
		}

		reqBody, _ := json.Marshal(req)

		res := testutil.NewRequest().
			Post("/posts").
			WithHeader("Content-Type", "application/json").
			WithBody(reqBody).
			WithJWSAuth(jwtToken).
			GoWithHTTPHandler(t, e)

		assert.Equal(t, http.StatusBadRequest, res.Code())

		var body server.RequestError
		err := res.UnmarshalBodyToObject(&body)
		assert.NoError(t, err)
		assert.Equal(t, errValidationFailed, body.Code)
		assert.Equal(t, "Post validation failed", body.Message)
	})
}

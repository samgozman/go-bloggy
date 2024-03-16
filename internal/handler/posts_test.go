package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/testutil"
	"github.com/samgozman/go-bloggy/internal/api"
	"github.com/samgozman/go-bloggy/internal/db"
	"github.com/samgozman/go-bloggy/internal/db/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"strconv"
	"testing"
	"time"
)

const jwtToken = "jwtToken"
const basePostsPath = "/posts"

func Test_PostPosts(t *testing.T) {
	conn, errDB := db.InitDatabase("file::memory:")
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

	t.Run("OK", func(t *testing.T) {
		e, _, mockJwtService, _, _ := registerHandlers(conn, []string{strconv.Itoa(user.ID)})
		mockJwtService.On("ParseTokenString", jwtToken).Return(user.ExternalID, nil)

		req := api.PostRequest{
			Title:       "Test Title",
			Slug:        uuid.New().String(),
			Content:     "Test Content to read in 1 second",
			Description: "Test Description",
			Keywords:    &[]string{"test1", "test2"},
		}

		reqBody, _ := json.Marshal(req)

		res := testutil.NewRequest().
			Post(basePostsPath).
			WithHeader("Content-Type", "application/json").
			WithBody(reqBody).
			WithJWSAuth(jwtToken).
			GoWithHTTPHandler(t, e)

		assert.Equal(t, http.StatusCreated, res.Code())

		var post api.PostResponse
		err := res.UnmarshalBodyToObject(&post)
		assert.NoError(t, err)

		assert.Equal(t, req.Title, post.Title)
		assert.Equal(t, req.Slug, post.Slug)
		assert.Equal(t, req.Content, post.Content)
		assert.Equal(t, req.Description, post.Description)
		assert.Equal(t, req.Keywords, post.Keywords)
		assert.Equal(t, 1, post.ReadingTime)
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
		e, _, mockJwtService, _, _ := registerHandlers(conn, nil)
		mockJwtService.On("ParseTokenString", jwtToken).Return(user.ExternalID, nil)

		req := api.PostRequest{
			Title:       "Test Title",
			Slug:        uuid.New().String(),
			Content:     "Test Content",
			Description: "Test Description",
			Keywords:    &[]string{"test1", "test2"},
		}

		reqBody, _ := json.Marshal(req)

		// Note: no Content-Type header
		res := testutil.NewRequest().
			Post(basePostsPath).
			WithBody(reqBody).
			WithJWSAuth(jwtToken).
			GoWithHTTPHandler(t, e)

		assert.Equal(t, http.StatusBadRequest, res.Code())

		var body api.RequestError
		err := res.UnmarshalBodyToObject(&body)
		assert.NoError(t, err)
		assert.Equal(t, errRequestBodyBinding, body.Code)
		assert.Equal(t, body.Message, fmt.Sprintf("Error binding request body: %v", echo.ErrUnsupportedMediaType.Message))
	})

	t.Run("400 - errGetUser", func(t *testing.T) {
		e, _, mockJwtService, _, _ := registerHandlers(conn, nil)
		mockJwtService.On("ParseTokenString", jwtToken).Return(uuid.New().String(), nil)

		req := api.PostRequest{
			Title:       "Test Title",
			Slug:        uuid.New().String(),
			Content:     "Test Content",
			Description: "Test Description",
		}

		reqBody, _ := json.Marshal(req)

		res := testutil.NewRequest().
			Post(basePostsPath).
			WithHeader("Content-Type", "application/json").
			WithBody(reqBody).
			WithJWSAuth(jwtToken).
			GoWithHTTPHandler(t, e)

		assert.Equal(t, http.StatusBadRequest, res.Code())

		var body api.RequestError
		err := res.UnmarshalBodyToObject(&body)
		assert.NoError(t, err)
		assert.Equal(t, errGetUser, body.Code)
		assert.Equal(t, "Post author is not found", body.Message)
	})

	t.Run("409 - errDuplicatePost", func(t *testing.T) {
		e, _, mockJwtService, _, _ := registerHandlers(conn, nil)
		mockJwtService.On("ParseTokenString", jwtToken).Return(user.ExternalID, nil)

		post1 := api.PostRequest{
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
			Post(basePostsPath).
			WithHeader("Content-Type", "application/json").
			WithBody(reqBody).
			WithJWSAuth(jwtToken).
			GoWithHTTPHandler(t, e)

		assert.Equal(t, http.StatusConflict, res.Code())

		var body api.RequestError
		err = res.UnmarshalBodyToObject(&body)
		assert.NoError(t, err)
		assert.Equal(t, errDuplicatePost, body.Code)
		assert.Equal(t, "Post with this URL slug already exists", body.Message)
	})

	t.Run("400 - errValidationFailed", func(t *testing.T) {
		e, _, mockJwtService, _, _ := registerHandlers(conn, nil)
		mockJwtService.On("ParseTokenString", jwtToken).Return(user.ExternalID, nil)

		req := api.PostRequest{
			Title:       "Test Title",
			Slug:        uuid.New().String(),
			Content:     "Test Content",
			Description: "Test Description",
			Keywords:    &[]string{"test1", "", ""}, // invalid keywords
		}

		reqBody, _ := json.Marshal(req)

		res := testutil.NewRequest().
			Post(basePostsPath).
			WithHeader("Content-Type", "application/json").
			WithBody(reqBody).
			WithJWSAuth(jwtToken).
			GoWithHTTPHandler(t, e)

		assert.Equal(t, http.StatusBadRequest, res.Code())

		var body api.RequestError
		err := res.UnmarshalBodyToObject(&body)
		assert.NoError(t, err)
		assert.Equal(t, errValidationFailed, body.Code)
		assert.Equal(t, "Post validation failed", body.Message)
	})
}

func TestHandler_GetPostsSlug(t *testing.T) {
	conn, errDB := db.InitDatabase("file::memory:")
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

	e, _, _, _, _ := registerHandlers(conn, []string{strconv.Itoa(user.ID)})

	// create post for test
	post := &models.Post{
		UserID:      user.ID,
		Title:       "Test Title",
		Slug:        "test-slug",
		Content:     "Test Content to read in 1 second",
		Description: "Test Description",
		Keywords:    "test1,test2",
	}
	err = conn.Models.Posts.Create(context.Background(), post)
	assert.NoError(t, err)

	t.Run("200 - OK", func(t *testing.T) {
		res := testutil.NewRequest().
			Get(basePostsPath+"/"+post.Slug).
			GoWithHTTPHandler(t, e)

		assert.Equal(t, http.StatusOK, res.Code())

		var postRes api.PostResponse
		err := res.UnmarshalBodyToObject(&postRes)
		assert.NoError(t, err)

		assert.Equal(t, post.Title, postRes.Title)
		assert.Equal(t, post.Slug, postRes.Slug)
		assert.Equal(t, post.Content, postRes.Content)
		assert.Equal(t, post.Description, postRes.Description)
		assert.Equal(t, &[]string{"test1", "test2"}, postRes.Keywords)
		assert.Equal(t, 1, postRes.ReadingTime)
		assert.NotEmpty(t, postRes.Id)
		assert.NotEmpty(t, postRes.CreatedAt)
		assert.NotEmpty(t, postRes.UpdatedAt)
	})

	t.Run("404 - Not Found", func(t *testing.T) {
		res := testutil.NewRequest().
			Get(basePostsPath+"/not-found-slug").
			GoWithHTTPHandler(t, e)

		assert.Equal(t, http.StatusNotFound, res.Code())

		var body api.RequestError
		err := res.UnmarshalBodyToObject(&body)
		assert.NoError(t, err)
		assert.Equal(t, errPostNotFound, body.Code)
		assert.Equal(t, "Post not found", body.Message)
	})

	t.Run("400 - errParamValidation", func(t *testing.T) {
		res := testutil.NewRequest().
			Get(basePostsPath+"/&kek*").
			GoWithHTTPHandler(t, e)

		assert.Equal(t, http.StatusBadRequest, res.Code())

		var body api.RequestError
		err := res.UnmarshalBodyToObject(&body)
		assert.NoError(t, err)
		assert.Equal(t, errParamValidation, body.Code)
		assert.Equal(t, "Slug is empty", body.Message)
	})
}

func TestHandler_GetPosts(t *testing.T) {
	conn, errDB := db.InitDatabase("file::memory:")
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

	e, _, _, _, _ := registerHandlers(conn, []string{strconv.Itoa(user.ID)})

	// create posts for test
	posts := []*models.Post{
		{
			UserID:      user.ID,
			Title:       "Test Title 1",
			Slug:        "test-slug-1",
			Content:     "Test Content 1",
			Description: "Test Description 1",
			Keywords:    "test1,test2",
		},
		{
			UserID:      user.ID,
			Title:       "Test Title 2",
			Slug:        "test-slug-2",
			Content:     "Test Content 1 to read in 1 second",
			Description: "Test Description 2",
			Keywords:    "test1,test2",
		},
	}
	for _, post := range posts {
		err = conn.Models.Posts.Create(context.Background(), post)
		assert.NoError(t, err)
	}

	t.Run("200 - OK", func(t *testing.T) {
		res := testutil.NewRequest().
			Get(basePostsPath).
			GoWithHTTPHandler(t, e)

		assert.Equal(t, http.StatusOK, res.Code())

		var postsRes api.PostsListResponse
		err := res.UnmarshalBodyToObject(&postsRes)
		assert.NoError(t, err)

		assert.Equal(t, postsRes.Total, 2)
		assert.Len(t, postsRes.Posts, 2)
	})

	t.Run("OK - with limit and offset", func(t *testing.T) {
		res := testutil.NewRequest().
			Get(basePostsPath+"?limit=1&page=1").
			GoWithHTTPHandler(t, e)

		assert.Equal(t, http.StatusOK, res.Code())

		var postsRes api.PostsListResponse
		err := res.UnmarshalBodyToObject(&postsRes)
		assert.NoError(t, err)

		assert.Equal(t, postsRes.Total, 2)
		assert.Len(t, postsRes.Posts, 1)
		assert.Equal(t, posts[1].Title, postsRes.Posts[0].Title)
		assert.Equal(t, posts[1].Slug, postsRes.Posts[0].Slug)
		assert.Equal(t, posts[1].Description, postsRes.Posts[0].Description)
		assert.Equal(t, &[]string{"test1", "test2"}, postsRes.Posts[0].Keywords)
		assert.Equal(t, 1, postsRes.Posts[0].ReadingTime)
	})

	t.Run("OK - zero posts per page", func(t *testing.T) {
		res := testutil.NewRequest().
			Get(basePostsPath+"?limit=1&page=100").
			GoWithHTTPHandler(t, e)

		assert.Equal(t, http.StatusOK, res.Code())

		var postsRes api.PostsListResponse
		err := res.UnmarshalBodyToObject(&postsRes)
		assert.NoError(t, err)

		assert.Equal(t, postsRes.Total, 2)
		assert.Len(t, postsRes.Posts, 0)
	})

	t.Run("400 - errParamValidation - limit", func(t *testing.T) {
		res := testutil.NewRequest().
			Get(basePostsPath+"?limit=0").
			GoWithHTTPHandler(t, e)

		assert.Equal(t, http.StatusBadRequest, res.Code())

		var body api.RequestError
		err := res.UnmarshalBodyToObject(&body)
		assert.NoError(t, err)
		assert.Equal(t, errParamValidation, body.Code)
		assert.Equal(t, "Limit must be between 1 and 25", body.Message)
	})

	t.Run("400 - errParamValidation - page", func(t *testing.T) {
		res := testutil.NewRequest().
			Get(basePostsPath+"?page=0").
			GoWithHTTPHandler(t, e)

		assert.Equal(t, http.StatusBadRequest, res.Code())

		var body api.RequestError
		err := res.UnmarshalBodyToObject(&body)
		assert.NoError(t, err)
		assert.Equal(t, errParamValidation, body.Code)
		assert.Equal(t, "Page must be greater than 0", body.Message)
	})
}

func TestHandler_PutPostsSlug(t *testing.T) {
	conn, errDB := db.InitDatabase("file::memory:")
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

	basePath := basePostsPath + "/"
	// create post for test
	post := &models.Post{
		UserID:      user.ID,
		Title:       "Test Title",
		Slug:        "test-slug",
		Content:     "Test Content to read in 1 second",
		Description: "Test Description",
		Keywords:    "test1,test2",
	}
	err = conn.Models.Posts.Create(context.Background(), post)
	assert.NoError(t, err)

	t.Run("OK", func(t *testing.T) {
		e, _, mockJwtService, _, _ := registerHandlers(conn, []string{strconv.Itoa(user.ID)})
		mockJwtService.On("ParseTokenString", jwtToken).Return(user.ExternalID, nil)

		req := api.PutPostRequest{
			Title:       "New Title",
			Content:     "New Content to read in 1 second",
			Description: "New Description",
			Keywords:    &[]string{"new1", "new2"},
		}

		reqBody, _ := json.Marshal(req)

		res := testutil.NewRequest().
			Put(basePath+post.Slug).
			WithHeader("Content-Type", "application/json").
			WithBody(reqBody).
			WithJWSAuth(jwtToken).
			GoWithHTTPHandler(t, e)

		assert.Equal(t, http.StatusOK, res.Code())

		var postRes api.PostResponse
		err := res.UnmarshalBodyToObject(&postRes)
		assert.NoError(t, err)

		assert.Equal(t, req.Title, postRes.Title)
		assert.Equal(t, post.Slug, postRes.Slug)
		assert.Equal(t, req.Content, postRes.Content)
		assert.Equal(t, req.Description, postRes.Description)
		assert.Equal(t, req.Keywords, postRes.Keywords)
		assert.Equal(t, 1, postRes.ReadingTime)
		assert.NotEmpty(t, postRes.Id)
		assert.NotEmpty(t, postRes.CreatedAt)
		assert.NotEmpty(t, postRes.UpdatedAt)

		// check that post is in the database
		postFromDB, err := conn.Models.Posts.GetBySlug(context.Background(), post.Slug)
		assert.NoError(t, err)
		assert.Equal(t, req.Title, postFromDB.Title)
		assert.Equal(t, post.Slug, postFromDB.Slug)
		assert.Equal(t, req.Content, postFromDB.Content)
		assert.Equal(t, req.Description, postFromDB.Description)
		assert.Equal(t, postFromDB.Keywords, "new1,new2")
		assert.NotEmpty(t, postFromDB.CreatedAt)
		assert.NotEmpty(t, postFromDB.UpdatedAt)
	})

	t.Run("400 - errRequestBodyBinding - ErrUnsupportedMediaType", func(t *testing.T) {
		e, _, mockJwtService, _, _ := registerHandlers(conn, nil)
		mockJwtService.On("ParseTokenString", jwtToken).Return(user.ExternalID, nil)

		req := api.PutPostRequest{
			Title:       "Test Title",
			Content:     "Test Content",
			Description: "Test Description",
			Keywords:    &[]string{"test1", "test2"},
		}

		reqBody, _ := json.Marshal(req)

		// Note: no Content-Type header
		res := testutil.NewRequest().
			Put(basePath+uuid.New().String()).
			WithBody(reqBody).
			WithJWSAuth(jwtToken).
			GoWithHTTPHandler(t, e)

		assert.Equal(t, http.StatusBadRequest, res.Code())

		var body api.RequestError
		err := res.UnmarshalBodyToObject(&body)
		assert.NoError(t, err)
		assert.Equal(t, errRequestBodyBinding, body.Code)
		assert.Equal(t, fmt.Sprintf("Error binding request body: %v", echo.ErrUnsupportedMediaType.Message), body.Message)
	})

	t.Run("404 - errPostNotFound", func(t *testing.T) {
		e, _, mockJwtService, _, _ := registerHandlers(conn, nil)
		mockJwtService.On("ParseTokenString", jwtToken).Return(user.ExternalID, nil)

		req := api.PutPostRequest{
			Title:       "Test Title",
			Content:     "Test Content",
			Description: "Test Description",
		}

		reqBody, _ := json.Marshal(req)

		res := testutil.NewRequest().
			Put("/posts/not-found-slug").
			WithHeader("Content-Type", "application/json").
			WithBody(reqBody).
			WithJWSAuth(jwtToken).
			GoWithHTTPHandler(t, e)

		assert.Equal(t, http.StatusNotFound, res.Code())

		var body api.RequestError
		err := res.UnmarshalBodyToObject(&body)
		assert.NoError(t, err)
		assert.Equal(t, errPostNotFound, body.Code)
		assert.Equal(t, "Post not found", body.Message)
	})

	t.Run("400 - errValidationFailed", func(t *testing.T) {
		e, _, mockJwtService, _, _ := registerHandlers(conn, []string{strconv.Itoa(user.ID)})
		mockJwtService.On("ParseTokenString", jwtToken).Return(user.ExternalID, nil)

		req := api.PutPostRequest{
			Title:       "New Title",
			Content:     "New Content to read in 1 second",
			Description: "New Description",
			Keywords:    &[]string{"new1", "", ""}, // invalid keywords
		}

		reqBody, _ := json.Marshal(req)

		res := testutil.NewRequest().
			Put(basePath+post.Slug).
			WithHeader("Content-Type", "application/json").
			WithBody(reqBody).
			WithJWSAuth(jwtToken).
			GoWithHTTPHandler(t, e)

		assert.Equal(t, http.StatusBadRequest, res.Code())

		var body api.RequestError
		err := res.UnmarshalBodyToObject(&body)
		assert.NoError(t, err)
		assert.Equal(t, errValidationFailed, body.Code)
		assert.Equal(t, "Post validation failed", body.Message)
	})
}

func TestHandler_PostPostsSlugSendEmail(t *testing.T) {
	conn, errDB := db.InitDatabase("file::memory:")
	if errDB != nil {
		t.Fatal(errDB)
	}

	// create user for test
	user := &models.User{
		ExternalID: uuid.New().String(),
		AuthMethod: models.GitHubAuthMethod,
		Login:      "testUser",
	}
	assert.NoError(t, conn.Models.Users.Upsert(context.Background(), user))

	// create post for test
	post := &models.Post{
		UserID:      user.ID,
		Title:       "Test Title",
		Slug:        uuid.New().String(),
		Content:     "Test Content to read in 1 second",
		Description: "Test Description",
		Keywords:    "test1,test2",
	}
	assert.NoError(t, conn.Models.Posts.Create(context.Background(), post))

	e, _, mockJwtService, mockMailerService, _ := registerHandlers(conn, []string{strconv.Itoa(user.ID)})
	mockJwtService.On("ParseTokenString", jwtToken).Return(user.ExternalID, nil)

	t.Run("201 - OK", func(t *testing.T) {
		// create subscription for test
		err := conn.Models.Subscribers.Create(context.Background(), &models.Subscriber{
			Email:       uuid.New().String() + "@test.com",
			IsConfirmed: true,
		})
		assert.NoError(t, err)

		mockMailerService.On("SendPostEmail", mock.Anything).Return(nil)

		res := testutil.NewRequest().
			Post(basePostsPath+"/"+post.Slug+"/send-email").
			WithJWSAuth(jwtToken).
			GoWithHTTPHandler(t, e)

		assert.Equal(t, http.StatusCreated, res.Code())
		mockMailerService.AssertExpectations(t)
	})

	t.Run("409 - errPostAlreadySent", func(t *testing.T) {
		// create a post for test that was already sent
		p := &models.Post{
			UserID:              user.ID,
			Title:               "Test Title",
			Slug:                uuid.New().String(),
			Content:             "Test Content to read in 1 second",
			Description:         "Test Description",
			SentToSubscribersAt: time.Now(),
		}
		assert.NoError(t, conn.Models.Posts.Create(context.Background(), p))

		res := testutil.NewRequest().
			Post(basePostsPath+"/"+p.Slug+"/send-email").
			WithJWSAuth(jwtToken).
			GoWithHTTPHandler(t, e)

		assert.Equal(t, http.StatusConflict, res.Code())

		var body api.RequestError
		err := res.UnmarshalBodyToObject(&body)
		assert.NoError(t, err)
		assert.Equal(t, errPostAlreadySent, body.Code)
		assert.Equal(t, "Post was already sent to subscribers. This can be done only once.", body.Message)
	})

	t.Run("404 - errPostNotFound", func(t *testing.T) {
		res := testutil.NewRequest().
			Post(basePostsPath+"/not-found-slug/send-email").
			WithJWSAuth(jwtToken).
			GoWithHTTPHandler(t, e)

		assert.Equal(t, http.StatusNotFound, res.Code())

		var body api.RequestError
		err := res.UnmarshalBodyToObject(&body)
		assert.NoError(t, err)
		assert.Equal(t, errPostNotFound, body.Code)
		assert.Equal(t, "Post not found", body.Message)
	})
}

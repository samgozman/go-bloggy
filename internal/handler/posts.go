package handler

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/samgozman/go-bloggy/internal/api"
	"github.com/samgozman/go-bloggy/internal/db/models"
	mailer "github.com/samgozman/go-bloggy/internal/mailer/types"
	"net/http"
	"regexp"
	"strings"
	"time"
)

func (h *Handler) PostPosts(ctx echo.Context) error {
	var req api.PostRequest
	if err := ctx.Bind(&req); err != nil {
		var errorMessage string
		var echoErr *echo.HTTPError
		if errors.As(err, &echoErr) {
			errorMessage = fmt.Sprintf("%v", echoErr.Message)
		}

		return ctx.JSON(http.StatusBadRequest, api.RequestError{
			Code:    errRequestBodyBinding,
			Message: fmt.Sprintf("Error binding request body: %v", errorMessage),
		})
	}

	var externalUserID string
	if s := ctx.Get("externalUserID"); s != nil {
		externalUserID = s.(string)
	}

	if externalUserID == "" {
		return ctx.JSON(http.StatusUnauthorized, api.RequestError{
			Code:    errUnauthorized,
			Message: "Unauthorized",
		})
	}

	user, err := h.db.Models().Users().GetByExternalID(ctx.Request().Context(), externalUserID)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, api.RequestError{
			Code:    errGetUser,
			Message: "Post author is not found",
		})
	}

	var keywords string
	if req.Keywords != nil && len(*req.Keywords) > 0 {
		keywords = (*req.Keywords)[0]
		for i := 1; i < len(*req.Keywords); i++ {
			keywords += "," + (*req.Keywords)[i]
		}
	}

	post := models.Post{
		UserID:      user.ID,
		Title:       req.Title,
		Slug:        req.Slug,
		Content:     req.Content,
		Description: req.Description,
		Keywords:    keywords,
	}

	if err := h.db.Models().Posts().Create(ctx.Request().Context(), &post); err != nil {
		switch {
		case errors.Is(err, models.ErrDuplicate):
			return ctx.JSON(http.StatusConflict, api.RequestError{
				Code:    errDuplicatePost,
				Message: "Post with this URL slug already exists",
			})
		case errors.Is(err, models.ErrValidationFailed):
			return ctx.JSON(http.StatusBadRequest, api.RequestError{
				Code:    errValidationFailed,
				Message: "Post validation failed",
			})

		default:
			return ctx.JSON(http.StatusInternalServerError, api.RequestError{
				Code:    errCreatePost,
				Message: "Error creating post",
			})
		}
	}

	return ctx.JSON(http.StatusCreated, api.PostResponse{
		Id:          post.ID,
		Title:       post.Title,
		Slug:        post.Slug,
		Description: post.Description,
		Content:     post.Content,
		Keywords:    req.Keywords,
		ReadingTime: post.ReadingTime,
		CreatedAt:   post.CreatedAt,
		UpdatedAt:   post.UpdatedAt,
	})
}

func (h *Handler) GetPostsSlug(ctx echo.Context, slug string) error {
	if !regexp.MustCompile(`^[a-z0-9-]+$`).MatchString(slug) {
		return ctx.JSON(http.StatusBadRequest, api.RequestError{
			Code:    errParamValidation,
			Message: "Slug is empty",
		})
	}

	post, err := h.db.Models().Posts().GetBySlug(ctx.Request().Context(), slug)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, api.RequestError{
			Code:    errPostNotFound,
			Message: "Post not found",
		})
	}

	var keywords []string
	if post.Keywords != "" {
		keywords = strings.Split(post.Keywords, ",")
	} else {
		keywords = []string{}
	}

	return ctx.JSON(http.StatusOK, api.PostResponse{
		Id:          post.ID,
		Title:       post.Title,
		Slug:        post.Slug,
		Description: post.Description,
		Content:     post.Content,
		Keywords:    &keywords,
		ReadingTime: post.ReadingTime,
		CreatedAt:   post.CreatedAt,
		UpdatedAt:   post.UpdatedAt,
	})
}

func (h *Handler) GetPosts(ctx echo.Context, params api.GetPostsParams) error {
	limit := 10
	page := 1
	if params.Limit != nil {
		limit = *params.Limit
	}

	if params.Page != nil {
		page = *params.Page
	}

	if limit < 1 || limit > 25 {
		return ctx.JSON(http.StatusBadRequest, api.RequestError{
			Code:    errParamValidation,
			Message: "Limit must be between 1 and 25",
		})
	}

	if page < 1 {
		return ctx.JSON(http.StatusBadRequest, api.RequestError{
			Code:    errParamValidation,
			Message: "Page must be greater than 0",
		})
	}

	count, err := h.db.Models().Posts().Count(ctx.Request().Context())
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.RequestError{
			Code:    errGetPostsCount,
			Message: "Error getting posts",
		})
	}

	posts, err := h.db.Models().Posts().FindAll(ctx.Request().Context(), page, limit)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.RequestError{
			Code:    errGetPosts,
			Message: "Error getting posts",
		})
	}

	postsItems := make([]api.PostsListItem, 0, len(posts))
	for _, post := range posts {
		var keywords []string
		if post.Keywords != "" {
			keywords = strings.Split(post.Keywords, ",")
		} else {
			keywords = []string{}
		}

		postsItems = append(postsItems, api.PostsListItem{
			Title:               post.Title,
			Slug:                post.Slug,
			Description:         post.Description,
			Keywords:            &keywords,
			ReadingTime:         post.ReadingTime,
			CreatedAt:           post.CreatedAt,
			SentToSubscribersAt: post.SentToSubscribersAt,
		})
	}

	return ctx.JSON(http.StatusOK, api.PostsListResponse{
		Posts: postsItems,
		Total: int(count),
	})
}

func (h *Handler) PutPostsSlug(ctx echo.Context, slug string) error {
	var req api.PutPostRequest
	if err := ctx.Bind(&req); err != nil {
		var errorMessage string
		var echoErr *echo.HTTPError
		if errors.As(err, &echoErr) {
			errorMessage = fmt.Sprintf("%v", echoErr.Message)
		}

		return ctx.JSON(http.StatusBadRequest, api.RequestError{
			Code:    errRequestBodyBinding,
			Message: fmt.Sprintf("Error binding request body: %v", errorMessage),
		})
	}

	post, err := h.db.Models().Posts().GetBySlug(ctx.Request().Context(), slug)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, api.RequestError{
			Code:    errPostNotFound,
			Message: "Post not found",
		})
	}

	post.Title = req.Title
	post.Description = req.Description
	post.Content = req.Content

	if req.Keywords != nil && len(*req.Keywords) > 0 {
		keywords := (*req.Keywords)[0]
		for i := 1; i < len(*req.Keywords); i++ {
			keywords += "," + (*req.Keywords)[i]
		}
		post.Keywords = keywords
	} else {
		post.Keywords = ""
	}

	if err := h.db.Models().Posts().Update(ctx.Request().Context(), post); err != nil {
		switch {
		case errors.Is(err, models.ErrDuplicate):
			return ctx.JSON(http.StatusConflict, api.RequestError{
				Code:    errDuplicatePost,
				Message: "Post with this URL slug already exists",
			})
		case errors.Is(err, models.ErrValidationFailed):
			return ctx.JSON(http.StatusBadRequest, api.RequestError{
				Code:    errValidationFailed,
				Message: "Post validation failed",
			})

		default:
			return ctx.JSON(http.StatusInternalServerError, api.RequestError{
				Code:    errUpdatePost,
				Message: "Error updating post",
			})
		}
	}

	return ctx.JSON(http.StatusOK, api.PostResponse{
		Id:          post.ID,
		Title:       post.Title,
		Slug:        post.Slug,
		Description: post.Description,
		Content:     post.Content,
		Keywords:    req.Keywords,
		ReadingTime: post.ReadingTime,
		CreatedAt:   post.CreatedAt,
		UpdatedAt:   post.UpdatedAt,
	})
}

func (h *Handler) PostPostsSlugSendEmail(ctx echo.Context, slug string) error {
	post, err := h.db.Models().Posts().GetBySlug(ctx.Request().Context(), slug)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, api.RequestError{
			Code:    errPostNotFound,
			Message: "Post not found",
		})
	}

	// check if post was already sent to subscribers
	if !post.SentToSubscribersAt.IsZero() {
		return ctx.JSON(http.StatusConflict, api.RequestError{
			Code:    errPostAlreadySent,
			Message: "Post was already sent to subscribers. This can be done only once.",
		})
	}

	// Get subscribers
	subs, err := h.db.Models().Subscribers().GetConfirmed(ctx.Request().Context())
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.RequestError{
			Code:    errGetSubscription,
			Message: "Error getting subscribers",
		})
	}

	if len(subs) == 0 {
		return ctx.JSON(http.StatusBadRequest, api.RequestError{
			Code:    errGetSubscription,
			Message: "No subscribers to send the post to.",
		})
	}

	// Transform subscribers
	mailerSubs := make([]*mailer.Subscriber, 0, len(subs))
	for _, sub := range subs {
		mailerSubs = append(mailerSubs, &mailer.Subscriber{
			Email: sub.Email,
			ID:    sub.ID.String(),
		})
	}

	// Send email to subscribers
	err = h.mailerService.SendPostEmail(&mailer.PostEmailSend{
		To:          mailerSubs,
		Title:       post.Title,
		Description: post.Description,
		Slug:        post.Slug,
	})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.RequestError{
			Code:    errSendPostEmail,
			Message: "Error sending post email",
		})
	}

	// Update post
	post.SentToSubscribersAt = time.Now()
	err = h.db.Models().Posts().Update(ctx.Request().Context(), post)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.RequestError{
			Code:    errUpdatePost,
			Message: "Error updating post",
		})
	}

	return ctx.NoContent(http.StatusCreated)
}

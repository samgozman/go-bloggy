package models

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"regexp"
	"strings"
	"time"
)

// AvgWordsPerMinute is the average number of words per minute a person can read.
const AvgWordsPerMinute = 250

// PostRepository is the database for the post data.
type PostRepository struct {
	conn *gorm.DB
}

// NewPostRepository creates a new PostRepository.
func NewPostRepository(conn *gorm.DB) *PostRepository {
	return &PostRepository{
		conn: conn,
	}
}

// Post is the model for the post-data.
type Post struct {
	ID                  int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Slug                string    `json:"slug" gorm:"uniqueIndex"` // Slug is the URL friendly version of the title
	Title               string    `json:"title"`
	Description         string    `json:"description"`
	Keywords            string    `json:"keywords"` // Keywords are comma separated
	Content             string    `json:"content"`
	ReadingTime         int       `json:"reading_time"` // ReadingTime is the estimated time to read the post in seconds
	UserID              int       `json:"user_id" gorm:"not null;constraint:OnUpdate:CASCADE;foreignKey:ID;references:ID"`
	SentToSubscribersAt time.Time `json:"sent_to_subscribers_at" gorm:"default:null"` // If not null, the post was sent
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

func (p *Post) Validate() error {
	switch {
	case p.Slug == "":
		return ErrPostURLRequired
	case !regexp.MustCompile(`^[a-z0-9-]+$`).MatchString(p.Slug):
		return ErrPostInvalidSlug
	case p.Title == "":
		return ErrPostTitleRequired
	case p.Description == "":
		return ErrPostDescriptionRequired
	case p.Content == "":
		return ErrPostContentRequired
	case p.Keywords != "":
		keywords := strings.Split(p.Keywords, ",")
		for _, k := range keywords {
			if k == "" {
				return ErrPostWrongKeywordsString
			}
		}
	case p.UserID == 0:
		return ErrPostUserIDRequired
	}

	return nil
}

func (p *Post) BeforeCreate(_ *gorm.DB) error {
	err := p.Validate()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrValidationFailed, err)
	}

	// Note: store the reading time in a database for a faster retrieval list of posts (without content)
	p.ReadingTime = int(p.CountReadingTime().Seconds())

	return nil
}

func (p *Post) BeforeUpdate(_ *gorm.DB) error {
	err := p.Validate()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrValidationFailed, err)
	}

	p.UpdatedAt = time.Now()
	p.ReadingTime = int(p.CountReadingTime().Seconds())

	return nil
}

// CountReadingTime counts the reading time of the post.
func (p *Post) CountReadingTime() time.Duration {
	words := strings.Fields(p.Content)
	readingTimeInMinutes := float64(len(words)) / AvgWordsPerMinute
	readingTimeInSeconds := readingTimeInMinutes * 60
	return time.Duration(readingTimeInSeconds) * time.Second
}

// PostRepositoryInterface is the interface for the PostRepository.
type PostRepositoryInterface interface {
	Create(ctx context.Context, p *Post) error
	GetBySlug(ctx context.Context, slug string) (*Post, error)
	FindAll(ctx context.Context, page, perPage int) ([]*Post, error)
	Update(ctx context.Context, p *Post) error
	Count(ctx context.Context) (int64, error)
}

// Create creates a new Post.
func (db *PostRepository) Create(ctx context.Context, p *Post) error {
	err := db.conn.WithContext(ctx).Create(p).Error
	if err != nil {
		return mapGormError(err)
	}

	return nil
}

// GetBySlug finds a Post by its URL Slug.
func (db *PostRepository) GetBySlug(ctx context.Context, slug string) (*Post, error) {
	var p Post
	err := db.conn.WithContext(ctx).Where("slug = ?", slug).First(&p).Error
	if err != nil {
		return nil, mapGormError(err)
	}

	return &p, nil
}

// FindAll returns all the posts with pagination, sorted by the created time.
// Selects only the necessary fields to reduce the payload - slug, title, description, keywords, created_at, sent_to_subscribers_at.
func (db *PostRepository) FindAll(ctx context.Context, page, perPage int) ([]*Post, error) {
	var posts []*Post
	err := db.conn.
		WithContext(ctx).
		Select("slug, title, description, keywords, reading_time, created_at, sent_to_subscribers_at").
		Offset((page - 1) * perPage).
		Limit(perPage).
		Order("created_at desc").
		Find(&posts).Error
	if err != nil {
		return nil, mapGormError(err)
	}

	return posts, nil
}

// Update updates the Post.
func (db *PostRepository) Update(ctx context.Context, p *Post) error {
	err := db.conn.WithContext(ctx).Save(p).Error
	if err != nil {
		return mapGormError(err)
	}

	return nil
}

// Count returns the total number of posts.
func (db *PostRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := db.conn.WithContext(ctx).Model(&Post{}).Count(&count).Error
	if err != nil {
		return 0, mapGormError(err)
	}

	return count, nil
}

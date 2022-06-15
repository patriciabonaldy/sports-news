package internal

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidData = errors.New("data can not be empty")

	ErrIDIsEmpty       = errors.New("invalid ID")
	ErrArticleNotFound = errors.New("id not found")
)

type Storage interface {
	GetByID(ctx context.Context, ID string) (ArticleNews, error)
	Save(ctx context.Context, news ArticleNews) error
}

//go:generate mockery --case=snake --outpkg=storagemocks --output=platform/storage/storagemocks --name=Storage

// ArticleNews is a structure of article to be stored
type ArticleNews struct {
	NewsID            string    `json:"news_id"`
	Title             string    `json:"title"`
	Subtitle          string    `json:"subtitle"`
	BodyText          string    `json:"body_text"`
	GalleryImageURLs  string    `json:"gallery_image_urls"`
	VideoURL          string    `json:"video_url"`
	Taxonomies        string    `json:"taxonomies"`
	TeaserText        string    `json:"teaser_text"`
	ThumbnailImageURL string    `json:"thumbnail_image_url"`
	PublishDate       string    `json:"publish_date"`
	IsPublished       bool      `json:"is_published"`
	CreateAt          time.Time `json:"create_at"`
}

func NewArticle() ArticleNews {
	id, _ := uuid.NewUUID()

	return ArticleNews{
		NewsID:      id.String(),
		IsPublished: false,
	}
}

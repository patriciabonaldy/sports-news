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
	NewsID            string
	Title             string
	Subtitle          string
	BodyText          string
	GalleryImageURLs  string
	VideoURL          string
	Taxonomies        string
	TeaserText        string
	ThumbnailImageURL string
	PublishDate       string
	IsPublished       bool
	CreateAt          time.Time
}

func NewArticle() ArticleNews {
	id, _ := uuid.NewUUID()

	return ArticleNews{
		NewsID:      id.String(),
		IsPublished: false,
	}
}

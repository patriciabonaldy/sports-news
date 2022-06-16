package business

import (
	"context"
	"github.com/patriciabonaldy/sports-news/internal"
	"github.com/patriciabonaldy/sports-news/internal/platform/logger"
)

type Service interface {
	GetArticles(ctx context.Context) ([]internal.ArticleNews, error)
	GetArticleByID(ctx context.Context, articleID string) (*internal.ArticleNews, error)
}

type service struct {
	repository internal.Storage
	log        logger.Logger
}

// NewService returns the default Service interface implementation.
func NewService(repository internal.Storage, log logger.Logger) Service {
	return &service{
		repository: repository,
		log:        log,
	}
}

func (s service) GetArticleByID(ctx context.Context, articleID string) (*internal.ArticleNews, error) {
	article, err := s.repository.GetArticleByID(ctx, articleID)
	if err != nil {
		s.log.Errorf("error GetArticleByID ID:%s:%s", articleID, err.Error())
		return nil, err
	}

	return &article, nil
}

func (s service) GetArticles(ctx context.Context) ([]internal.ArticleNews, error) {
	articles, err := s.repository.GetArticles(ctx)
	if err != nil {
		s.log.Errorf("error Get Articles:%s", err.Error())
		return nil, err
	}

	return articles, nil
}

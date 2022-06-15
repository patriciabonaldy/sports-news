package providers

import (
	"context"
	"encoding/json"
	"time"

	"github.com/patriciabonaldy/sports-news/internal"
	"github.com/patriciabonaldy/sports-news/internal/platform/logger"
	"github.com/patriciabonaldy/sports-news/internal/platform/pubsub"
)

type service struct {
	subscriber pubsub.Subscriber
	repository internal.Storage
	log        logger.Logger
}

func NewBrenfordSubscriber(subscriber pubsub.Subscriber, repository internal.Storage, log logger.Logger) *service {
	return &service{
		subscriber: subscriber,
		repository: repository,
		log:        log,
	}
}

func (s *service) Start(ctx context.Context) error {
	return s.subscriber.Subscriber(ctx, s.callBack)
}

func (s *service) callBack(ctx context.Context, message interface{}) error {
	data, err := json.Marshal(message)
	if err != nil {
		s.log.Error("invalid message type")
		return err
	}

	var msg pubsub.Message
	err = json.Unmarshal(data, &msg)
	if err != nil {
		s.log.Error("invalid message type")
		return err
	}

	var msgSync messageSync
	err = json.Unmarshal(msg.RawData, &msgSync)
	if err != nil {
		s.log.Error("invalid message type")
		return err
	}

	article := toArticle(msgSync)
	err = s.repository.Save(ctx, article)
	if err != nil {
		s.log.Errorf("error CreateArticleNews %s", err.Error())
		return err
	}

	return nil
}

func toArticle(msg messageSync) internal.ArticleNews {
	article := internal.ArticleNews{
		NewsID:            msg.NewsArticleID,
		Title:             msg.Title,
		Subtitle:          "",
		BodyText:          "",
		GalleryImageURLs:  "",
		VideoURL:          "",
		Taxonomies:        msg.Taxonomies,
		TeaserText:        msg.TeaserText,
		ThumbnailImageURL: msg.ThumbnailImageURL,
		PublishDate:       msg.PublishDate,
		IsPublished:       msg.IsPublished,
		CreateAt:          time.Now(),
	}

	return article
}
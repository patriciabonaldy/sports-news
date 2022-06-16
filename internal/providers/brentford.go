package providers

import (
	"context"
	"encoding/json"

	"github.com/patriciabonaldy/sports-news/internal/platform/logger"
	"github.com/patriciabonaldy/sports-news/internal/platform/pubsub"
)

type service struct {
	pipeline   Pipeline
	subscriber pubsub.Subscriber
	log        logger.Logger
}

func NewBrenfordSubscriber(pipeline Pipeline, subscriber pubsub.Subscriber, log logger.Logger) *service {
	return &service{
		pipeline:   pipeline,
		subscriber: subscriber,
		log:        log,
	}
}

func (s *service) Start(ctx context.Context) error {
	return s.subscriber.Subscriber(ctx, s.callBack)
}

func (s *service) callBack(ctx context.Context, message interface{}) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	var msg pubsub.Message
	err = json.Unmarshal(data, &msg)
	if err != nil {
		return err
	}

	var msgSync []NewsletterNewsItem
	err = json.Unmarshal(msg.RawData, &msgSync)
	if err != nil {
		return err
	}

	s.pipeline.Process(ctx, msgSync)
	s.log.Info("Brenford sync process finished")

	return nil
}

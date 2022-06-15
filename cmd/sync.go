package main

import (
	"context"
	"github.com/patriciabonaldy/big_queue/pkg/kafka"
	"github.com/patriciabonaldy/sports-news/internal/platform/pubsub"
	"github.com/patriciabonaldy/sports-news/internal/platform/syncer/news"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"strings"

	"github.com/patriciabonaldy/sports-news/internal/config"
	"github.com/patriciabonaldy/sports-news/internal/platform/logger"
)

const (
	newsSyncCronSpec = "1 * * * *"
)

func sync(cfg *config.Config, cron *cron.Cron, log logger.Logger) error {
	log.Info("sync", "provider news")
	if _, err := cron.AddFunc(newsSyncCronSpec, providerNews(cfg, log)); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func providerNews(cfg *config.Config, log logger.Logger) func() {
	return func() {
		log.Info("synchronizing")
		if cfg.Kafka.Topic == "" {
			log.Info("topic-id was not configured")
			return
		}

		ctx := context.Background()
		publisher := kafka.NewPublisher(strings.Split(cfg.Kafka.Broker, ","), cfg.Kafka.Topic)
		producer := pubsub.NewProducer(publisher)

		s := news.NewSyncer(nil, producer, log, cfg)
		if err := s.Sync(ctx); err != nil {
			log.Error(err)
		}
	}
}

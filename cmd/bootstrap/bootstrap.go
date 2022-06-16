package bootstrap

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/robfig/cron/v3"

	"github.com/patriciabonaldy/big_queue/pkg/kafka"
	"github.com/patriciabonaldy/sports-news/internal/business"
	"github.com/patriciabonaldy/sports-news/internal/config"
	"github.com/patriciabonaldy/sports-news/internal/platform/genericClient"
	"github.com/patriciabonaldy/sports-news/internal/platform/logger"
	"github.com/patriciabonaldy/sports-news/internal/platform/pubsub"
	"github.com/patriciabonaldy/sports-news/internal/platform/server"
	"github.com/patriciabonaldy/sports-news/internal/platform/server/handler"
	"github.com/patriciabonaldy/sports-news/internal/platform/storage/mongo"
	"github.com/patriciabonaldy/sports-news/internal/providers"
)

const (
	lsbnTZ = "Europe/Lisbon"
)

func Run() error {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	logger := logger.New()
	loc, err := time.LoadLocation(lsbnTZ)
	if err != nil {
		loc = time.Local
	}

	c := cron.New(cron.WithLocation(loc))
	if err = sync(cfg, c, logger); err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	repository, err := mongo.NewDBStorage(ctx, cfg.Database, logger)
	if err != nil {
		log.Fatal(err)
	}

	svc := business.NewService(repository, logger)
	handler := handler.New(svc, logger)
	ctx, srv := server.New(ctx, cfg, handler)

	runBrenfordSubscriber(ctx, cfg, repository, logger)
	c.Start()

	return srv.Run(ctx)
}

func runBrenfordSubscriber(ctx context.Context, cfg *config.Config, repository *mongo.Repository, log logger.Logger) error {
	if cfg.Kafka.Topic == "" {
		log.Info("topic-id was not configured")
		return errors.New("topic-id was not configured")
	}

	client := genericClient.New()
	pipeline := providers.NewPipeLine(repository, client, log)
	consumer := kafka.NewConsumer(strings.Split(cfg.Kafka.Broker, ","), cfg.Kafka.Topic)
	subscriber := pubsub.NewSubscriber(consumer, log)
	pSubscriber := providers.NewBrenfordSubscriber(pipeline, subscriber, log)

	go pSubscriber.Start(ctx)

	return nil
}

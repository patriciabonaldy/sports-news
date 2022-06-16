package pubsub

import (
	"context"
	"log"

	"github.com/patriciabonaldy/big_queue/pkg"
	"github.com/patriciabonaldy/sports-news/internal/platform/logger"
)

type subscriber struct {
	consumer pkg.Consumer
	log      logger.Logger
}

type Subscriber interface {
	Subscriber(ctx context.Context, callback func(ctx context.Context, message interface{}) error)
}

//go:generate mockery --case=snake --outpkg=pubsubMock --output=pubsubMock --name=Subscriber

func NewSubscriber(consumer pkg.Consumer, log logger.Logger) Subscriber {
	p := subscriber{
		consumer: consumer,
		log:      log,
	}

	return &p
}

func (s subscriber) Subscriber(ctx context.Context, callback func(ctx context.Context, message interface{}) error) {
	chMsg := make(chan pkg.Message)
	chErr := make(chan error)
	go func() {
		s.consumer.Read(ctx, chMsg, chErr)
	}()

	// read/process message
	for {
		select {
		case m := <-chMsg:
			err := callback(ctx, m)
			if err != nil {
				log.Println(err)
			}
		case err := <-chErr:
			log.Println(err)
		}
	}
}

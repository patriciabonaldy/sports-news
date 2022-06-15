package sync_news

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"io"

	"github.com/patriciabonaldy/sports-news/internal/platform/genericClient"
	"github.com/patriciabonaldy/sports-news/internal/platform/logger"
	"github.com/patriciabonaldy/sports-news/internal/platform/pubsub"
	"github.com/patriciabonaldy/sports-news/internal/platform/syncer"
)

type SyncerNews struct {
	log      logger.Logger
	client   genericClient.Client
	producer pubsub.Producer
	url      string
}

var _ syncer.Syncer = &SyncerNews{}

func NewSyncerNews(log logger.Logger, client genericClient.Client, url string) syncer.Syncer {
	return &SyncerNews{log: log, client: client, url: url}
}

func (s *SyncerNews) Sync(ctx context.Context) error {
	resp, err := s.client.Get(ctx, s.url)
	if err != nil {
		s.log.Errorf("error fetch - sync %s", err)
		return err
	}

	defer resp.Body.Close()
	var newListInf newListInformation
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		s.log.Errorf("error read body - sync %s", err)
		return err
	}

	err = xml.Unmarshal(data, &newListInf)
	if err != nil {
		s.log.Errorf("error unmarshall - sync %s", err)
		return err
	}

	m, err := generateMessage(newListInf.NewsletterNews.Article)
	if err != nil {
		s.log.Errorf("error generate message - sync %s", err)
		return err
	}

	err = s.producer.Produce(ctx, m)
	if err != nil {
		s.log.Errorf("error producer - sync %s", err)
		return err
	}

	return nil
}

func generateMessage(articles []article) (*pubsub.Message, error) {
	message, err := pubsub.NewSystemMessage()
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(articles)
	if err != nil {
		return nil, err
	}

	message.RawData = data

	return &message, nil
}

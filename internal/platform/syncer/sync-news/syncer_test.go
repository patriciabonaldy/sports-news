package sync_news

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/patriciabonaldy/sports-news/internal/platform/genericClient"
	"github.com/patriciabonaldy/sports-news/internal/platform/pubsub"
	"github.com/patriciabonaldy/sports-news/internal/platform/pubsub/pubsubMock"
)

func TestSyncerNews_Sync(t *testing.T) {
	type fields struct {
		client         genericClient.Client
		fnMockProducer func() pubsub.Producer
		url            string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Error in get",
			fields: fields{
				client:         &mockClient{wantError: true},
				fnMockProducer: func() pubsub.Producer { return nil },
				url:            "",
			},
			wantErr: true,
		},
		{
			name: "Error unmarshall body response",
			fields: fields{
				client:         &mockClient{wantErrorUnmarshall: true},
				fnMockProducer: func() pubsub.Producer { return nil },
				url:            "",
			},
			wantErr: true,
		},
		{
			name: "Error publish message",
			fields: fields{
				client: &mockClient{},
				fnMockProducer: func() pubsub.Producer {
					productMock := new(pubsubMock.Producer)
					productMock.On("Produce", mock.Anything, mock.Anything).
						Return(errors.New("something unexpected happened"))

					return productMock
				},
				url: "",
			},
			wantErr: true,
		},
		{
			name: "success",
			fields: fields{
				client: &mockClient{},
				fnMockProducer: func() pubsub.Producer {
					productMock := new(pubsubMock.Producer)
					productMock.On("Produce", mock.Anything, mock.Anything).
						Return(nil)

					return productMock
				},
				url: "",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SyncerNews{
				log:      &mockLog{},
				client:   tt.fields.client,
				producer: tt.fields.fnMockProducer(),
				url:      tt.fields.url,
			}
			if err := s.Sync(context.Background()); (err != nil) != tt.wantErr {
				t.Errorf("Sync() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

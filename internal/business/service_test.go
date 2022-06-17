package business

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"

	"github.com/patriciabonaldy/sports-news/internal"
	"github.com/patriciabonaldy/sports-news/internal/platform/logger"
	"github.com/patriciabonaldy/sports-news/internal/platform/storage/storagemocks"
)

func Test_service_GetArticleByID(t *testing.T) {
	tests := []struct {
		name    string
		repo    func() internal.Storage
		want    func() *internal.ArticleNews
		wantErr bool
	}{
		{
			name: "error getting article",
			repo: func() internal.Storage {
				repoMock := new(storagemocks.Storage)
				repoMock.On("GetArticleByID", mock.Anything, mock.Anything).
					Return(&internal.ArticleNews{}, errors.New("something unexpected happened"))

				return repoMock

			},
			want: func() *internal.ArticleNews {
				return nil
			},
			wantErr: true,
		},
		{
			name: "success",
			repo: func() internal.Storage {
				mockA := mockArticle()
				repoMock := new(storagemocks.Storage)
				repoMock.On("GetArticleByID", mock.Anything, mock.Anything).
					Return(&mockA, nil)

				return repoMock

			},
			want: func() *internal.ArticleNews {
				mockA := mockArticle()

				return &mockA
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService(tt.repo(), logger.New())
			got, err := s.GetArticleByID(context.Background(), "641838")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetArticleByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			want := tt.want()
			if !reflect.DeepEqual(got, want) {
				t.Errorf("GetArticleByID() got = %v, want %v", got, want)
			}
		})
	}
}

func Test_service_GetArticles(t *testing.T) {
	tests := []struct {
		name    string
		repo    func() internal.Storage
		want    []internal.ArticleNews
		wantErr bool
	}{
		{
			name: "error getting article",
			repo: func() internal.Storage {
				repoMock := new(storagemocks.Storage)
				repoMock.On("GetArticles", mock.Anything).
					Return(nil, errors.New("something unexpected happened"))

				return repoMock

			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "success",
			repo: func() internal.Storage {
				repoMock := new(storagemocks.Storage)
				repoMock.On("GetArticles", mock.Anything).
					Return([]internal.ArticleNews{mockArticle()}, nil)

				return repoMock

			},
			want: []internal.ArticleNews{mockArticle()},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService(tt.repo(), logger.New())
			got, err := s.GetArticles(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("GetArticles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetArticles() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func mockArticle() internal.ArticleNews {
	return internal.ArticleNews{
		NewsID:         "641838",
		ClubName:       "Brentford",
		ClubWebsiteURL: "https://www.brentfordfc.com",
		ArticleURL:     "https://www.brentfordfc.com/news/2022/june/pontus-explains-how-fatherhood-has-calmed-him-down",
		Title:          "Pontus explains",
		Subtitle:       "Pontus explains ",
		BodyText:       "Pontus explains Pontus explains Pontus explains",
		IsPublished:    true,
		CreateAt:       time.Time{},
	}
}

package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/patriciabonaldy/sports-news/internal"
	"github.com/patriciabonaldy/sports-news/internal/business"
	"github.com/patriciabonaldy/sports-news/internal/platform/logger"
	"github.com/patriciabonaldy/sports-news/internal/platform/storage/storagemocks"
)

var timeN = time.Now()

func TestHandler_GetArticleByID(t *testing.T) {
	repositoryMock := new(storagemocks.Storage)
	repositoryMock.On("GetArticleByID", mock.Anything, mock.Anything).
		Return(&internal.ArticleNews{}, errors.New("something unexpected happened")).Once()

	articleMock := mockArticle()
	repositoryMock.On("GetArticleByID", mock.Anything, mock.Anything).
		Return(&articleMock, nil).Once()
	log := logger.New()
	svc := business.NewService(repositoryMock, log)
	handler := New(svc, log)
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/articles/:id", handler.GetArticleByID())

	t.Run("given a invalid request it returns 400", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/articles/0", nil)
		require.NoError(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})

	t.Run("given a error it returns 500", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/articles/8001122", nil)
		require.NoError(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	})

	t.Run("given a valid request it returns 200", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/articles/8001122", nil)
		require.NoError(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		var resp Response
		err = json.NewDecoder(res.Body).Decode(&resp)
		require.NoError(t, err)

		want := Response{
			NewsID:         "641838",
			ClubName:       "Brentford",
			ClubWebsiteURL: "https://www.brentfordfc.com",
			ArticleURL:     "https://www.brentfordfc.com/news/2022/june/pontus-explains-how-fatherhood-has-calmed-him-down",
			Title:          "Pontus explains",
			Subtitle:       "Pontus explains ",
			BodyText:       "Pontus explains Pontus explains Pontus explains",
			IsPublished:    true,
			CreateAt:       timeN,
		}
		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, resp.NewsID, want.NewsID)
	})
}

func TestHandler_GetArticles(t *testing.T) {
	repositoryMock := new(storagemocks.Storage)
	repositoryMock.On("GetArticles", mock.Anything).
		Return(nil, errors.New("something unexpected happened")).Once()

	repositoryMock.On("GetArticles", mock.Anything).
		Return([]internal.ArticleNews{mockArticle()}, nil).Once()
	log := logger.New()
	svc := business.NewService(repositoryMock, log)
	handler := New(svc, log)
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/articles", handler.GetArticles())

	t.Run("given a error it returns 500", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/articles", nil)
		require.NoError(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	})

	t.Run("given a valid request it returns 200", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/articles", nil)
		require.NoError(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		var resp []Response
		err = json.NewDecoder(res.Body).Decode(&resp)
		require.NoError(t, err)

		want := []Response{
			{
				NewsID:         "641838",
				ClubName:       "Brentford",
				ClubWebsiteURL: "https://www.brentfordfc.com",
				ArticleURL:     "https://www.brentfordfc.com/news/2022/june/pontus-explains-how-fatherhood-has-calmed-him-down",
				Title:          "Pontus explains",
				Subtitle:       "Pontus explains ",
				BodyText:       "Pontus explains Pontus explains Pontus explains",
				IsPublished:    true,
				CreateAt:       timeN,
			},
		}
		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, len(resp), len(want))
	})
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
		CreateAt:       timeN,
	}
}

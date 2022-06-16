package providers

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"strconv"
	"sync"
	"time"

	"github.com/patriciabonaldy/sports-news/internal"
	"github.com/patriciabonaldy/sports-news/internal/platform/genericClient"
	"github.com/patriciabonaldy/sports-news/internal/platform/logger"
)

type Pipeline interface {
	Process(ctx context.Context, data []NewsletterNewsItem)
}

type pipeLine struct {
	repository internal.Storage
	client     genericClient.Client
	log        logger.Logger
}

func NewPipeLine(repository internal.Storage, client genericClient.Client, log logger.Logger) Pipeline {
	return &pipeLine{repository: repository, client: client, log: log}
}

func (p *pipeLine) Process(ctx context.Context, data []NewsletterNewsItem) {
	ch1 := p.taskFetch(ctx, data)
	ch2 := p.taskParse(ch1)

	var wg sync.WaitGroup
	wg.Add(len(data))
	go func(ch2 chan NewsArticleInformation) {
		for m := range ch2 {
			article := toArticle(m)
			_, err := p.repository.GetArticleByID(ctx, article.NewsID)
			if err != nil {
				err = p.repository.Save(ctx, article)
				if err != nil {
					p.log.Errorf("error Create ArticleNews %s", err.Error())
				}
			}

			wg.Done()
		}
	}(ch2)

	wg.Wait()
}

func (p *pipeLine) taskFetch(ctx context.Context, data []NewsletterNewsItem) chan []byte {
	ch1 := make(chan []byte)
	for _, d := range data {
		go func(ch chan []byte, id string) {
			resp, err := p.client.Get(ctx, fmt.Sprintf("https://www.brentfordfc.com/api/incrowd/getnewsarticleinformation?id=%s", id))
			if err != nil {
				p.log.Errorf("error fetch - sync articleID %s, error %s", id, err)
				return
			}

			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				p.log.Errorf("error read body - sync %s", err)
				return
			}

			ch <- body
		}(ch1, d.NewsArticleID)

	}

	return ch1
}

func (p *pipeLine) taskParse(ch1 chan []byte) chan NewsArticleInformation {
	ch2 := make(chan NewsArticleInformation)

	go func(ch1 chan []byte, ch2 chan NewsArticleInformation) {
		for data := range ch1 {
			var newArticle NewsArticleInformation
			err := xml.Unmarshal(data, &newArticle)
			if err != nil {
				p.log.Errorf("error unmarshall - sync %s", err)
				return
			}

			ch2 <- newArticle
		}
	}(ch1, ch2)

	return ch2
}

func toArticle(msg NewsArticleInformation) internal.ArticleNews {
	newsArticle := msg.NewsArticle
	text := ""
	if data, err := json.Marshal(newsArticle.BodyText); err == nil {
		text = string(data)
	}

	isPublished, err := strconv.ParseBool(newsArticle.IsPublished)
	if err != nil {
		isPublished = false
	}

	article := internal.ArticleNews{
		NewsID:            newsArticle.NewsArticleID,
		Title:             newsArticle.Title,
		Subtitle:          newsArticle.Subtitle,
		ClubName:          msg.ClubName,
		ClubWebsiteURL:    msg.ClubWebsiteURL,
		ArticleURL:        newsArticle.ArticleURL,
		BodyText:          text,
		GalleryImageURLs:  newsArticle.GalleryImageURLs,
		VideoURL:          newsArticle.VideoURL,
		Taxonomies:        newsArticle.Taxonomies,
		TeaserText:        newsArticle.TeaserText,
		ThumbnailImageURL: newsArticle.ThumbnailImageURL,
		PublishDate:       newsArticle.PublishDate,
		IsPublished:       isPublished,
		CreateAt:          time.Now(),
	}

	return article
}

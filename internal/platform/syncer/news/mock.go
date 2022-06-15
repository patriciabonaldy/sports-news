package news

import (
	"context"
	"errors"
	"github.com/patriciabonaldy/sports-news/internal/platform/logger"
	"io"
	"net/http"
	"strings"

	"github.com/patriciabonaldy/sports-news/internal/platform/genericClient"
)

type mockClient struct {
	wantError           bool
	wantErrorUnmarshall bool
}

func (m mockClient) Delete(_ context.Context, _ string, _ ...genericClient.Header) error {
	if m.wantError {
		return errors.New("unknown error")
	}

	return nil
}

func (m mockClient) Get(_ context.Context, _ string) (resp *http.Response, err error) {
	if m.wantError {
		return nil, errors.New("unknown error")
	}

	reader := io.NopCloser(strings.NewReader(`<NewListInformation>
<ClubName>Brentford</ClubName>
<ClubWebsiteURL>https://www.brentfordfc.com</ClubWebsiteURL>
<NewsletterNewsItems>
	<NewsletterNewsItem>
		<ArticleURL>https://www.brentfordfc.com/news/2022/june/international-round-up-15.06.22/</ArticleURL>
		<NewsArticleID>641772</NewsArticleID>
		<PublishDate>2022-06-15 10:00:00</PublishDate>
		<Taxonomies>Players</Taxonomies>
		<TeaserText></TeaserText>
		<ThumbnailImageURL>https://www.brentfordfc.com/api/image/feedassets/1fa93314-18a6-4180-bd1e-6e7a3549c969/Medium/mads-bidstrup-denmark-u21.jpg</ThumbnailImageURL>
		<Title>Three wins for international Bees yesterday</Title>
		<OptaMatchId></OptaMatchId>
		<LastUpdateDate>2022-06-15 09:53:56</LastUpdateDate>
		<IsPublished>True</IsPublished>
	</NewsletterNewsItem>
	<NewsletterNewsItem>
		<ArticleURL>https://www.brentfordfc.com/news/2022/june/202122---brentfords-fourth-highest-league-finish/</ArticleURL>
		<NewsArticleID>641745</NewsArticleID>
		<PublishDate>2022-06-15 08:00:00</PublishDate>
		<Taxonomies>History</Taxonomies>
		<TeaserText></TeaserText>
		<ThumbnailImageURL>https://www.brentfordfc.com/api/image/feedassets/377baa5d-ea74-41dc-9527-e73a8692c07c/Medium/dai-hopkins-brentford-v-portsmouth-1939.jpg</ThumbnailImageURL>
		<Title>2021/22 - Brentford&apos;s fourth-highest league finish</Title>
		<OptaMatchId></OptaMatchId>
		<LastUpdateDate>2022-06-15 08:00:21</LastUpdateDate>
		<IsPublished>True</IsPublished>
		</NewsletterNewsItem>
	</NewsletterNewsItems>
</NewListInformation>
		`))
	if m.wantErrorUnmarshall {
		reader = io.NopCloser(strings.NewReader(`<NewListInformation`))
	}

	return &http.Response{
		StatusCode:    http.StatusCreated,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		ContentLength: 0,
		Body:          reader,
	}, nil
}

func (m mockClient) Post(_ context.Context, _ string, _ []byte, _ ...genericClient.Header) (resp *http.Response, err error) {
	if m.wantError {
		return nil, errors.New("unknown error")
	}

	reader := io.NopCloser(strings.NewReader(`{
		"data": {}
		`))

	return &http.Response{
		StatusCode:    http.StatusCreated,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		ContentLength: 0,
		Body:          reader,
	}, nil
}

var _ genericClient.Client = &mockClient{}

type mockLog struct{}

func (m mockLog) Error(args ...interface{}) {
}

func (m mockLog) Errorf(format string, args ...interface{}) {

}

func (m mockLog) Info(args ...interface{}) {

}

func (m mockLog) Infof(format string, args ...interface{}) {

}

var _ logger.Logger = &mockLog{}

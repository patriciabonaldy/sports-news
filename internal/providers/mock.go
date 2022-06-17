package providers

import (
	"context"
	"errors"
	"github.com/patriciabonaldy/sports-news/internal/platform/genericClient"
	"io"
	"net/http"
	"strings"
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
		<ArticleURL>https://www.brentfordfc.com/news/2022/june/202122---brentfords-fourth-highest-league-finish/</ArticleURL>
		<NewsID>641745</NewsID>
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

func mockNewsletterNewsItem() NewsletterNewsItem {
	return NewsletterNewsItem{
		ArticleURL:    "https://www.brentfordfc.com/news/2022/june/pontus-explains-how-fatherhood-has-calmed-him-down",
		NewsArticleID: "641838",
		Title:         "Pontus explains",
		IsPublished:   "true",
	}
}

package providers

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/patriciabonaldy/sports-news/internal"
	"github.com/patriciabonaldy/sports-news/internal/platform/genericClient"
	"github.com/patriciabonaldy/sports-news/internal/platform/logger"
)

func Test_pipeLine_taskFetch(t *testing.T) {
	type fields struct {
		repository internal.Storage
		client     genericClient.Client
	}
	tests := []struct {
		name       string
		fields     fields
		expectData bool
		want       []byte
	}{
		{
			name: "want error fetch data",
			fields: fields{
				client: &mockClient{wantError: true},
			},
		},
		{
			name: "success",
			fields: fields{
				client: &mockClient{},
			},
			expectData: true,
			want: []byte(`<NewListInformation>
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
		`),
		},
	}
	data := []NewsletterNewsItem{mockNewsletterNewsItem()}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &pipeLine{
				repository: tt.fields.repository,
				client:     tt.fields.client,
				log:        logger.New(),
			}

			ch := p.taskFetch(context.Background(), data)

			var got string
			if tt.expectData {
				got = string(<-ch)
			}

			want := string(tt.want)
			if !reflect.DeepEqual(got, want) {
				t.Errorf("taskFetch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_pipeLine_taskParse(t *testing.T) {
	tests := []struct {
		name       string
		ch1        func() chan []byte
		want       func() chan NewsArticleInformation
		expectData bool
	}{
		{
			name: "want error unmarshal",
			ch1: func() chan []byte {
				ch2 := make(chan []byte)
				go func() {
					ch2 <- []byte(`<>`)
				}()
				return ch2
			},
			want: func() chan NewsArticleInformation {
				return make(chan NewsArticleInformation)
			},
		},
		{
			name: "success",
			ch1: func() chan []byte {
				ch2 := make(chan []byte)
				go func() {
					ch2 <- []byte(`<NewsArticleInformation>
									<ClubName>Brentford</ClubName>
									<ClubWebsiteURL>https://www.brentfordfc.com</ClubWebsiteURL>
									<NewsArticle>
									<ArticleURL>https://www.brentfordfc.com/news/2017/june/be-there-in-201718/</ArticleURL>
									<NewsArticleID>173860</NewsArticleID>
									<PublishDate>2017-06-05 10:33:39</PublishDate>
									<Taxonomies>Ticket News</Taxonomies>
									<TeaserText/>
									<Subtitle/>
									<ThumbnailImageURL/>
									<Title>Be There for our 2017/18 season</Title>
									<BodyText/>
									<GalleryImageURLs/>
									<VideoURL/>
									<OptaMatchId/>
									<LastUpdateDate>2019-09-02 03:36:54</LastUpdateDate>
									<IsPublished>True</IsPublished>
									</NewsArticle>
									</NewsArticleInformation>
											`)
				}()
				return ch2
			},
			want: func() chan NewsArticleInformation {
				ch2 := make(chan NewsArticleInformation)
				go func() {
					ch2 <- NewsArticleInformation{
						ClubName:       "Brentford",
						ClubWebsiteURL: "https://www.brentfordfc.com",
						NewsArticle: NewsArticle{
							Text:           "\n\t\t\t\t\t\t\t\t\t\n\t\t\t\t\t\t\t\t\t\n\t\t\t\t\t\t\t\t\t\n\t\t\t\t\t\t\t\t\t\n\t\t\t\t\t\t\t\t\t\n\t\t\t\t\t\t\t\t\t\n\t\t\t\t\t\t\t\t\t\n\t\t\t\t\t\t\t\t\t\n\t\t\t\t\t\t\t\t\t\n\t\t\t\t\t\t\t\t\t\n\t\t\t\t\t\t\t\t\t\n\t\t\t\t\t\t\t\t\t\n\t\t\t\t\t\t\t\t\t\n\t\t\t\t\t\t\t\t\t\n\t\t\t\t\t\t\t\t\t",
							ArticleURL:     "https://www.brentfordfc.com/news/2017/june/be-there-in-201718/",
							NewsArticleID:  "173860",
							PublishDate:    "2017-06-05 10:33:39",
							Taxonomies:     "Ticket News",
							Title:          "Be There for our 2017/18 season",
							LastUpdateDate: "2019-09-02 03:36:54",
							IsPublished:    "True",
						},
					}
				}()
				return ch2
			},
			expectData: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &pipeLine{
				log: logger.New(),
			}
			ch := p.taskParse(tt.ch1())
			chWant := tt.want()

			var got string
			var want string
			if tt.expectData {
				aiGot := <-ch
				bGot, err := json.Marshal(aiGot.NewsArticle)
				assert.NoError(t, err)
				got = string(bGot)

				aiWant := <-chWant
				bWant, err := json.Marshal(aiWant.NewsArticle)
				assert.NoError(t, err)
				want = string(bWant)
			}

			if !reflect.DeepEqual(got, want) {
				t.Errorf("taskParse() got= %v\n, want %v", got, want)
			}
		})
	}
}

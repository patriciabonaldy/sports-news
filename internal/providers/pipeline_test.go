package providers

import (
	"context"
	"reflect"
	"testing"

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
	type fields struct {
		repository internal.Storage
		client     genericClient.Client
		log        logger.Logger
	}
	type args struct {
		ch1 chan []byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   chan NewsArticleInformation
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &pipeLine{
				repository: tt.fields.repository,
				client:     tt.fields.client,
				log:        tt.fields.log,
			}
			if got := p.taskParse(tt.args.ch1); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("taskParse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_toArticle(t *testing.T) {
	type args struct {
		msg NewsArticleInformation
	}
	tests := []struct {
		name string
		args args
		want internal.ArticleNews
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toArticle(tt.args.msg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toArticle() = %v, want %v", got, tt.want)
			}
		})
	}
}

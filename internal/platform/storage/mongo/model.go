package mongo

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	_ "gopkg.in/mgo.v2/bson"

	"github.com/patriciabonaldy/sports-news/internal"
)

// ArticleNews is a structure of article to be stored
type ArticleNews struct {
	ArticleID         string              `bson:"article_id"`
	Title             string              `bson:"title"`
	Subtitle          string              `bson:"subtitle,omitempty"`
	BodyText          string              `bson:"body_text,omitempty"`
	GalleryImageURLs  string              `bson:"gallery_image_urls,omitempty"`
	VideoURL          string              `bson:"video_url,omitempty"`
	Taxonomies        string              `bson:"taxonomies,omitempty"`
	TeaserText        string              `bson:"teaser_text,omitempty"`
	ThumbnailImageURL string              `bson:"thumbnail_image_url,omitempty"`
	PublishDate       string              `bson:"publish_date,omitempty"`
	IsPublished       bool                `bson:"is_published,omitempty"`
	CreateAt          primitive.Timestamp `bson:"create_at"`
}

func (a *ArticleNews) createAt() time.Time {
	return time.Unix(int64(a.CreateAt.T), 0).UTC()
}

func parseToBusinessArticleNews(result ArticleNews) internal.ArticleNews {
	article := internal.ArticleNews{
		NewsID:   result.ArticleID,
		CreateAt: result.createAt(),
	}

	return article
}

func parseToArticleNewsDB(result internal.ArticleNews) ArticleNews {
	article := ArticleNews{
		ArticleID: result.NewsID,
		CreateAt: primitive.Timestamp{
			T: uint32(result.CreateAt.Unix()),
		},
	}

	return article
}

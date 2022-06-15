package brentfordFC

import "encoding/xml"

type newListInformation struct {
	XMLName        xml.Name       `xml:"NewListInformation"`
	Text           string         `xml:",chardata"`
	ClubName       string         `xml:"ClubName"`
	ClubWebsiteURL string         `xml:"ClubWebsiteURL"`
	NewsletterNews newsletterNews `xml:"NewsletterNewsItems"`
}

type newsletterNews struct {
	Text    string    `xml:",chardata"`
	Article []article `xml:"NewsletterNewsItem"`
}

type article struct {
	Text              string `xml:",chardata"`
	ArticleURL        string `xml:"article_url"`
	NewsArticleID     string `xml:"news_article_id"`
	PublishDate       string `xml:"publish_date"`
	Taxonomies        string `xml:"taxonomies"`
	TeaserText        string `xml:"teaser_text"`
	ThumbnailImageURL string `xml:"thumbnail_image_url"`
	Title             string `xml:"title"`
	OptaMatchId       string `xml:"opta_match_id"`
	LastUpdateDate    string `xml:"last_update_date"`
	IsPublished       string `xml:"is_published"`
}

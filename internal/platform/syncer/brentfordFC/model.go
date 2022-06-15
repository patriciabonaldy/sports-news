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
	ArticleURL        string `xml:"ArticleURL"`
	NewsArticleID     string `xml:"NewsID"`
	PublishDate       string `xml:"PublishDate"`
	Taxonomies        string `xml:"Taxonomies"`
	TeaserText        string `xml:"TeaserText"`
	ThumbnailImageURL string `xml:"ThumbnailImageURL"`
	Title             string `xml:"Title"`
	OptaMatchId       string `xml:"OptaMatchId"`
	LastUpdateDate    string `xml:"LastUpdateDate"`
	IsPublished       string `xml:"IsPublished"`
}

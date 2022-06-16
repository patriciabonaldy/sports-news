package brentfordFC

import "encoding/xml"

type NewListInformation struct {
	XMLName             xml.Name `xml:"NewListInformation"`
	Text                string   `xml:",chardata"`
	ClubName            string   `xml:"ClubName"`
	ClubWebsiteURL      string   `xml:"ClubWebsiteURL"`
	NewsletterNewsItems struct {
		Text               string               `xml:",chardata"`
		NewsletterNewsItem []NewsletterNewsItem `xml:"NewsletterNewsItem"`
	} `xml:"NewsletterNewsItems"`
}

type NewsletterNewsItem struct {
	Text              string `xml:",chardata"`
	ArticleURL        string `xml:"ArticleURL"`
	NewsArticleID     string `xml:"NewsArticleID"`
	PublishDate       string `xml:"PublishDate"`
	Taxonomies        string `xml:"Taxonomies"`
	TeaserText        string `xml:"TeaserText"`
	ThumbnailImageURL string `xml:"ThumbnailImageURL"`
	Title             string `xml:"Title"`
	OptaMatchId       string `xml:"OptaMatchId"`
	LastUpdateDate    string `xml:"LastUpdateDate"`
	IsPublished       string `xml:"IsPublished"`
}

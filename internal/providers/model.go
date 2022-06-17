package providers

import "encoding/xml"

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

type NewsArticle struct {
	Text              string `xml:",chardata"`
	ArticleURL        string `xml:"ArticleURL"`
	NewsArticleID     string `xml:"NewsArticleID"`
	PublishDate       string `xml:"PublishDate"`
	Taxonomies        string `xml:"Taxonomies"`
	TeaserText        string `xml:"TeaserText"`
	Subtitle          string `xml:"Subtitle"`
	ThumbnailImageURL string `xml:"ThumbnailImageURL"`
	Title             string `xml:"Title"`
	BodyText          struct {
		Text string `xml:",chardata"`
		P    []struct {
			Text string `xml:",chardata"`
			A    struct {
				Text string `xml:",chardata"`
				Href string `xml:"href,attr"`
			} `xml:"a"`
		} `xml:"p"`
	} `xml:"BodyText"`
	GalleryImageURLs string `xml:"GalleryImageURLs"`
	VideoURL         string `xml:"VideoURL"`
	OptaMatchId      string `xml:"OptaMatchId"`
	LastUpdateDate   string `xml:"LastUpdateDate"`
	IsPublished      string `xml:"IsPublished"`
}

type NewsArticleInformation struct {
	XMLName        xml.Name    `xml:"NewsArticleInformation"`
	Text           string      `xml:",chardata"`
	ClubName       string      `xml:"ClubName"`
	ClubWebsiteURL string      `xml:"ClubWebsiteURL"`
	NewsArticle    NewsArticle `xml:"NewsArticle"`
}

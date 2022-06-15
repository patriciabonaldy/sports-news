package providers

type messageSync struct {
	Text              string `json:",chardata"`
	ArticleURL        string `json:"article_url"`
	NewsArticleID     string `json:"news_article_id"`
	PublishDate       string `json:"publish_date"`
	Taxonomies        string `json:"taxonomies"`
	TeaserText        string `json:"teaser_text"`
	ThumbnailImageURL string `json:"thumbnail_image_url"`
	Title             string `json:"title"`
	OptaMatchId       string `json:"opta_match_id"`
	LastUpdateDate    string `json:"last_update_date"`
	IsPublished       bool   `json:"is_published"`
}

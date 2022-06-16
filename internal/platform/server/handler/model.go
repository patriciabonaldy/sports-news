package handler

import (
	"time"
)

// swagger:model RequestID
type RequestID struct {
	ID string `uri:"id" binding:"required" example:"8001122"`
}

// swagger:model Response
type Response struct {
	NewsID            string    `json:"news_id"`
	ClubName          string    `json:"club_name"`
	ClubWebsiteURL    string    `json:"club_website_url"`
	ArticleURL        string    `json:"article_url"`
	Title             string    `json:"title"`
	Subtitle          string    `json:"subtitle"`
	BodyText          string    `json:"body_text"`
	GalleryImageURLs  string    `json:"gallery_image_urls"`
	VideoURL          string    `json:"video_url"`
	Taxonomies        string    `json:"taxonomies"`
	TeaserText        string    `json:"teaser_text"`
	ThumbnailImageURL string    `json:"thumbnail_image_url"`
	PublishDate       string    `json:"publish_date"`
	IsPublished       bool      `json:"is_published"`
	CreateAt          time.Time `json:"create_at"`
}

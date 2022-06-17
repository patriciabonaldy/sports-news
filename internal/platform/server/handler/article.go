package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/patriciabonaldy/sports-news/internal"
	"github.com/patriciabonaldy/sports-news/internal/business"
	"github.com/patriciabonaldy/sports-news/internal/platform/logger"
)

type ArticleHandler struct {
	service business.Service
	log     logger.Logger
}

func New(service business.Service, log logger.Logger) ArticleHandler {
	return ArticleHandler{
		service: service,
		log:     log,
	}
}

func (a *ArticleHandler) GetArticles() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ans, err := a.service.GetArticles(ctx)
		if err != nil {
			switch err {
			case internal.ErrIDIsEmpty,
				internal.ErrInvalidData:
				ctx.JSON(http.StatusBadRequest, err.Error())
				return

			default:
				ctx.JSON(http.StatusInternalServerError, err.Error())
				return
			}
		}

		ctx.JSON(http.StatusOK, toResponseArticles(ans))
	}
}

func (a *ArticleHandler) GetArticleByID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req RequestID
		if err := ctx.ShouldBindUri(&req); err != nil {
			ctx.JSON(400, gin.H{"msg": err.Error()})
			return
		}

		if req.ID == "0" {
			ctx.JSON(400, gin.H{"msg": "bad request"})
			return
		}

		ans, err := a.service.GetArticleByID(ctx, req.ID)
		if err != nil {
			switch err {
			case internal.ErrIDIsEmpty,
				internal.ErrArticleNotFound,
				internal.ErrInvalidData:
				ctx.JSON(http.StatusBadRequest, err.Error())
				return

			default:
				ctx.JSON(http.StatusInternalServerError, err.Error())
				return
			}
		}

		ctx.JSON(http.StatusOK, toResponse(ans))
	}
}

func toResponseArticles(articleNews []internal.ArticleNews) []Response {
	var resp = make([]Response, 0, 1)
	for _, a := range articleNews {
		resp = append(resp, toResponse(&a))
	}

	return resp
}

func toResponse(articleNews *internal.ArticleNews) Response {
	return Response{
		NewsID:            articleNews.NewsID,
		ClubName:          articleNews.ClubName,
		ClubWebsiteURL:    articleNews.ClubWebsiteURL,
		ArticleURL:        articleNews.ArticleURL,
		Title:             articleNews.Title,
		Subtitle:          articleNews.Subtitle,
		BodyText:          articleNews.BodyText,
		GalleryImageURLs:  articleNews.GalleryImageURLs,
		VideoURL:          articleNews.VideoURL,
		Taxonomies:        articleNews.Taxonomies,
		TeaserText:        articleNews.TeaserText,
		ThumbnailImageURL: articleNews.ThumbnailImageURL,
		PublishDate:       articleNews.PublishDate,
		IsPublished:       articleNews.IsPublished,
		CreateAt:          articleNews.CreateAt,
	}
}

package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CheckHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "pong")
	}
}

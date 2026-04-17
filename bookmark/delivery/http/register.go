package http

import (
	"github.com/gin-gonic/gin"
	"github.com/zhashkevych/go-clean-architecture/bookmark"
)

func RegisterHTTPEndpoints(router *gin.RouterGroup, uc bookmark.UseCase) {
	h := NewHandler(uc)

	bookmarks := router.Group("/bookmarks")
	{
		bookmarks.POST("", h.Create)
		bookmarks.GET("", h.Get)
		bookmarks.GET("/by-tags", h.GetByTags)
		bookmarks.DELETE("", h.Delete)
		bookmarks.PUT("/tags", h.UpdateTags)
	}

	tags := router.Group("/tags")
	{
		tags.GET("", h.GetAllTags)
		tags.POST("/merge", h.MergeTags)
		tags.POST("/batch-add", h.BatchAddTags)
		tags.POST("/batch-remove", h.BatchRemoveTags)
	}
}

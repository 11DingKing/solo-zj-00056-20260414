package http

import (
	"github.com/gin-gonic/gin"
	"github.com/zhashkevych/go-clean-architecture/auth"
	"github.com/zhashkevych/go-clean-architecture/bookmark"
	"github.com/zhashkevych/go-clean-architecture/models"
	"net/http"
	"strings"
)

type Bookmark struct {
	ID    string   `json:"id"`
	URL   string   `json:"url"`
	Title string   `json:"title"`
	Tags  []string `json:"tags"`
}

type Handler struct {
	useCase bookmark.UseCase
}

func NewHandler(useCase bookmark.UseCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

type createInput struct {
	URL   string   `json:"url"`
	Title string   `json:"title"`
	Tags  []string `json:"tags"`
}

func (h *Handler) Create(c *gin.Context) {
	inp := new(createInput)
	if err := c.BindJSON(inp); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user := c.MustGet(auth.CtxUserKey).(*models.User)

	if err := h.useCase.CreateBookmark(c.Request.Context(), user, inp.URL, inp.Title, inp.Tags); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

type getResponse struct {
	Bookmarks []*Bookmark `json:"bookmarks"`
}

func (h *Handler) Get(c *gin.Context) {
	user := c.MustGet(auth.CtxUserKey).(*models.User)

	bms, err := h.useCase.GetBookmarks(c.Request.Context(), user)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, &getResponse{
		Bookmarks: toBookmarks(bms),
	})
}

func (h *Handler) GetByTags(c *gin.Context) {
	user := c.MustGet(auth.CtxUserKey).(*models.User)

	tagsParam := c.Query("tags")
	if tagsParam == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	tags := strings.Split(tagsParam, ",")

	bms, err := h.useCase.GetBookmarksByTags(c.Request.Context(), user, tags)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, &getResponse{
		Bookmarks: toBookmarks(bms),
	})
}

type deleteInput struct {
	ID string `json:"id"`
}

func (h *Handler) Delete(c *gin.Context) {
	inp := new(deleteInput)
	if err := c.BindJSON(inp); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user := c.MustGet(auth.CtxUserKey).(*models.User)

	if err := h.useCase.DeleteBookmark(c.Request.Context(), user, inp.ID); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

type updateTagsInput struct {
	ID   string   `json:"id"`
	Tags []string `json:"tags"`
}

func (h *Handler) UpdateTags(c *gin.Context) {
	inp := new(updateTagsInput)
	if err := c.BindJSON(inp); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user := c.MustGet(auth.CtxUserKey).(*models.User)

	if err := h.useCase.UpdateBookmarkTags(c.Request.Context(), user, inp.ID, inp.Tags); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

type mergeTagsInput struct {
	FromTag string `json:"from_tag"`
	ToTag   string `json:"to_tag"`
}

func (h *Handler) MergeTags(c *gin.Context) {
	inp := new(mergeTagsInput)
	if err := c.BindJSON(inp); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user := c.MustGet(auth.CtxUserKey).(*models.User)

	if err := h.useCase.MergeTags(c.Request.Context(), user, inp.FromTag, inp.ToTag); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

type batchTagsInput struct {
	BookmarkIDs []string `json:"bookmark_ids"`
	Tags        []string `json:"tags"`
}

func (h *Handler) BatchAddTags(c *gin.Context) {
	inp := new(batchTagsInput)
	if err := c.BindJSON(inp); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user := c.MustGet(auth.CtxUserKey).(*models.User)

	if err := h.useCase.BatchAddTags(c.Request.Context(), user, inp.BookmarkIDs, inp.Tags); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) BatchRemoveTags(c *gin.Context) {
	inp := new(batchTagsInput)
	if err := c.BindJSON(inp); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user := c.MustGet(auth.CtxUserKey).(*models.User)

	if err := h.useCase.BatchRemoveTags(c.Request.Context(), user, inp.BookmarkIDs, inp.Tags); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

type getTagsResponse struct {
	Tags []string `json:"tags"`
}

func (h *Handler) GetAllTags(c *gin.Context) {
	user := c.MustGet(auth.CtxUserKey).(*models.User)

	tags, err := h.useCase.GetAllTags(c.Request.Context(), user)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, &getTagsResponse{
		Tags: tags,
	})
}

func toBookmarks(bs []*models.Bookmark) []*Bookmark {
	out := make([]*Bookmark, len(bs))

	for i, b := range bs {
		out[i] = toBookmark(b)
	}

	return out
}

func toBookmark(b *models.Bookmark) *Bookmark {
	return &Bookmark{
		ID:    b.ID,
		URL:   b.URL,
		Title: b.Title,
		Tags:  b.Tags,
	}
}

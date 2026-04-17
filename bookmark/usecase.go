package bookmark

import (
	"context"
	"github.com/zhashkevych/go-clean-architecture/models"
)

type UseCase interface {
	CreateBookmark(ctx context.Context, user *models.User, url, title string, tags []string) error
	GetBookmarks(ctx context.Context, user *models.User) ([]*models.Bookmark, error)
	GetBookmarksByTags(ctx context.Context, user *models.User, tags []string) ([]*models.Bookmark, error)
	DeleteBookmark(ctx context.Context, user *models.User, id string) error
	UpdateBookmarkTags(ctx context.Context, user *models.User, id string, tags []string) error
	MergeTags(ctx context.Context, user *models.User, fromTag string, toTag string) error
	BatchAddTags(ctx context.Context, user *models.User, bookmarkIDs []string, tags []string) error
	BatchRemoveTags(ctx context.Context, user *models.User, bookmarkIDs []string, tags []string) error
	GetAllTags(ctx context.Context, user *models.User) ([]string, error)
}

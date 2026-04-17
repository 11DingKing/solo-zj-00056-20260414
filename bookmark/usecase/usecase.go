package usecase

import (
	"context"
	"github.com/zhashkevych/go-clean-architecture/bookmark"
	"github.com/zhashkevych/go-clean-architecture/models"
)

type BookmarkUseCase struct {
	bookmarkRepo bookmark.Repository
}

func NewBookmarkUseCase(bookmarkRepo bookmark.Repository) *BookmarkUseCase {
	return &BookmarkUseCase{
		bookmarkRepo: bookmarkRepo,
	}
}

func (b BookmarkUseCase) CreateBookmark(ctx context.Context, user *models.User, url, title string, tags []string) error {
	bm := &models.Bookmark{
		URL:   url,
		Title: title,
		Tags:  tags,
	}

	return b.bookmarkRepo.CreateBookmark(ctx, user, bm)
}

func (b BookmarkUseCase) GetBookmarks(ctx context.Context, user *models.User) ([]*models.Bookmark, error) {
	return b.bookmarkRepo.GetBookmarks(ctx, user)
}

func (b BookmarkUseCase) GetBookmarksByTags(ctx context.Context, user *models.User, tags []string) ([]*models.Bookmark, error) {
	return b.bookmarkRepo.GetBookmarksByTags(ctx, user, tags)
}

func (b BookmarkUseCase) DeleteBookmark(ctx context.Context, user *models.User, id string) error {
	return b.bookmarkRepo.DeleteBookmark(ctx, user, id)
}

func (b BookmarkUseCase) UpdateBookmarkTags(ctx context.Context, user *models.User, id string, tags []string) error {
	return b.bookmarkRepo.UpdateBookmarkTags(ctx, user, id, tags)
}

func (b BookmarkUseCase) MergeTags(ctx context.Context, user *models.User, fromTag string, toTag string) error {
	return b.bookmarkRepo.MergeTags(ctx, user, fromTag, toTag)
}

func (b BookmarkUseCase) BatchAddTags(ctx context.Context, user *models.User, bookmarkIDs []string, tags []string) error {
	return b.bookmarkRepo.BatchAddTags(ctx, user, bookmarkIDs, tags)
}

func (b BookmarkUseCase) BatchRemoveTags(ctx context.Context, user *models.User, bookmarkIDs []string, tags []string) error {
	return b.bookmarkRepo.BatchRemoveTags(ctx, user, bookmarkIDs, tags)
}

func (b BookmarkUseCase) GetAllTags(ctx context.Context, user *models.User) ([]string, error) {
	return b.bookmarkRepo.GetAllTags(ctx, user)
}

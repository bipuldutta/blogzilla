package usecases

import (
	"context"

	"github.com/bipuldutta/blogzilla/domain"
	"github.com/bipuldutta/blogzilla/utils"
)

var blogLogger = *utils.Logger()

/*
BlogManager is the actual business logic section for managing all blog related transactions
while this is a skeleton and just making calls to the repo layer at this time there could be
actual BL we could implement at some point
*/
type BlogManager struct {
	blogRepo domain.BlogRepo
}

func NewBlogManager(blogRepo domain.BlogRepo) *BlogManager {
	return &BlogManager{blogRepo: blogRepo}
}

func (m *BlogManager) Create(ctx context.Context, newBlog *domain.Blog) (int64, error) {
	// TODO figure out what to validate about the blog data
	return m.blogRepo.Create(ctx, newBlog)
}

func (m *BlogManager) Get(ctx context.Context, blogID int64) (*domain.Blog, error) {
	return m.blogRepo.Get(ctx, blogID)
}

func (m *BlogManager) Search(ctx context.Context, offset int, limit int, search string) ([]*domain.Blog, error) {
	blogLogger.Infof("offset: %d, limit: %d, search: %s", offset, limit, search)
	return m.blogRepo.Search(ctx, offset, limit, search)
}

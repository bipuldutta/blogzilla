package repositories

import (
	"context"
	"fmt"

	"github.com/bipuldutta/blogzilla/config"
	"github.com/bipuldutta/blogzilla/domain"
	"github.com/bipuldutta/blogzilla/utils"

	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	createBlogQuery = `INSERT INTO blogs (user_id, title, content, tags) VALUES ($1, $2, $3, $4) RETURNING id`
	searchBlogQuery = `
		SELECT id, user_id, title, content, tags, created_at, updated_at FROM blogs 
		WHERE (COALESCE(title, '') || ' ' || COALESCE(content, '') || ' ' || COALESCE(tags, '')) ILIKE '%' || $1 || '%'
		ORDER BY created_at DESC 
		OFFSET $2 LIMIT $3
    `
)

var blogLogger = *utils.Logger()

type BlogRepo struct {
	conf   *config.Config
	client *pgxpool.Pool
}

func NewBlogRepo(conf *config.Config, client *pgxpool.Pool) domain.BlogRepo {
	return &BlogRepo{
		conf:   conf,
		client: client,
	}
}

func (r *BlogRepo) Create(ctx context.Context, newBlog *domain.Blog) (int64, error) {
	var blogID int64
	// create blog and return its id
	err := r.client.QueryRow(ctx, createBlogQuery, newBlog.UserID, newBlog.Title, newBlog.Content, newBlog.Tags).Scan(&blogID)
	if err != nil {
		blogLogger.WithError(err).Error("failed to create blog")
		return -1, err
	}
	return blogID, nil
}

func (r *BlogRepo) Get(ctx context.Context, blogID int64) (*domain.Blog, error) {
	return nil, fmt.Errorf("not yet implemented")
}

func (r *BlogRepo) Search(ctx context.Context, offset int, limit int, search string) ([]*domain.Blog, error) {
	var blogs []*domain.Blog

	rows, err := r.client.Query(ctx, searchBlogQuery, search, offset, limit)
	if err != nil {
		blogLogger.WithError(err).Errorf("failed to query blogs. offset: %d, limit: %d, search: %s", offset, limit, search)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var blog domain.Blog
		err = rows.Scan(&blog.ID, &blog.UserID, &blog.Title, &blog.Content, &blog.Tags, &blog.CreatedAt, &blog.UpdatedAt)
		if err != nil {
			blogLogger.WithError(err).Errorf("failed to query blogs. offset: %d, limit: %d, search: %s", offset, limit, search)
			return nil, err
		}
		blogs = append(blogs, &blog)
	}

	return blogs, nil
}

package api

import (
	"time"
)

type CreateUserRequestV1 struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type UpdateUserRequestV1 struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type UserResponseV1 struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CreateBlogRequestV1 struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Tags    string `json:"tags"`
}

type UpdateBlogRequestV1 struct {
	ID      int64    `json:"id"`
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

type BlogResponseV1 struct {
	ID        int64          `json:"id"`
	Creator   *BlogCreatorV1 `json:"creator"`
	Title     string         `json:"title"`
	Content   string         `json:"content"`
	Tags      string         `json:"tags"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
}

type BlogCreatorV1 struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type CredentialsRequestV1 struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

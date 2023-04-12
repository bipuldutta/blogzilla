package domain

import "context"

type UserRepo interface {
	Create(ctx context.Context, user *User) (int64, error)
	Login(ctx context.Context, username string, password string) (string, error)
}

type AuthRepo interface {
	GetToken(ctx context.Context, userID int64, permissions []string) (string, error)
}

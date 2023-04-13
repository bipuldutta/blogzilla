package domain

import "context"

type DatabaseRepo interface {
	Initialize(ctx context.Context) error
}

type UserRepo interface {
	Create(ctx context.Context, user *User) (int64, error)
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	GetRoleByName(ctx context.Context, roleName string) (*Role, error)
	AssignRoles(ctx context.Context, userID int64, roleIDs ...int64) error
	Login(ctx context.Context, username string, password string) (string, error)
}

type AuthRepo interface {
	GetToken(ctx context.Context, userID int64, permissions []string) (string, error)
}

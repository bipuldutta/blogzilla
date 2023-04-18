package domain

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

/*
	These are the domain object for this application
*/

type User struct {
	ID        int64
	Username  string
	Password  string
	FirstName string
	LastName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Role struct {
	ID          int64
	Name        string
	Description string
	Permissions []string
}

type Blog struct {
	ID        int64
	UserID    int64
	Title     string
	Content   string
	Tags      string // comma separated
	CreatedAt time.Time
	UpdatedAt time.Time
}

// CustomClaims represents the custom claims for the JWT token.
type CustomClaims struct {
	UserID      int64
	Permissions map[string]any
	jwt.RegisteredClaims
}

func (cc *CustomClaims) HasPermission(permission string) bool {
	if _, ok := cc.Permissions[permission]; ok {
		return true
	}
	return false
}

package domain

import (
	"time"

	"github.com/dgrijalva/jwt-go"
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
	Permissions []string
	jwt.StandardClaims
}

func (cc *CustomClaims) HasPermission(permission string) bool {
	// TODO improve this so that the lookup is not O(n), maybe a map?
	for _, perm := range cc.Permissions {
		if permission == perm {
			return true
		}
	}
	return false
}

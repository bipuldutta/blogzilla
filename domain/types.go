package domain

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

/*
	TODO: Ideally we will not use the domain objects in the API layer
	because we probably don't want to send all the data back to the client
	due to time constraint I am using the same object
*/

type User struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type Blog struct {
	ID        int       `json:"id"`
	UserID    int       `json:"userId"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// CustomClaims represents the custom claims for the JWT token.
type CustomClaims struct {
	UserID      int64    `json:"user_id"`
	Permissions []string `json:"permissions"`
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

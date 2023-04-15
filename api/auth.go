package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/bipuldutta/blogzilla/config"
	"github.com/bipuldutta/blogzilla/usecases"
)

type AuthMiddleware struct {
	conf        *config.Config
	authManager *usecases.AuthManager
}

func NewAuthMiddleware(conf *config.Config, authManager *usecases.AuthManager) *AuthMiddleware {
	return &AuthMiddleware{
		conf:        conf,
		authManager: authManager,
	}
}

func (am *AuthMiddleware) authorize(permission string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := am.extractTokenFromHeader(r)
		if err != nil {
			http.Error(w, "invalid auth token", http.StatusUnauthorized)
			return
		}

		userID, err := am.authManager.ValidateToken(token, permission)
		if err != nil {
			http.Error(w, "invalid auth token", http.StatusUnauthorized)
			return
		}

		// inject the userId so that it can be collected downstream off the context
		ctx := context.WithValue(r.Context(), "userId", userID)
		r = r.WithContext(ctx)

		// Call next handler function in chain
		next.ServeHTTP(w, r)
	})
}

// extractTokenFromHeader extracts the JWT token from the Authorization header in the format "Bearer {token}".
func (am *AuthMiddleware) extractTokenFromHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("missing authorization header")
	}

	bearerToken := strings.Split(authHeader, " ")
	if len(bearerToken) != 2 || strings.ToLower(bearerToken[0]) != "bearer" {
		return "", fmt.Errorf("invalid authorization header format")
	}

	return bearerToken[1], nil
}

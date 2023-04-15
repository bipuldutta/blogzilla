package usecases

import (
	"fmt"
	"time"

	"github.com/bipuldutta/blogzilla/config"
	"github.com/bipuldutta/blogzilla/domain"
	"github.com/dgrijalva/jwt-go"
)

type AuthManager struct {
	conf *config.Config
}

func NewAuthManager(conf *config.Config) *AuthManager {
	return &AuthManager{
		conf: conf,
	}
}

func (m *AuthManager) ValidateToken(tokenString string, permission string) (int64, error) {
	// Parse the token without verifying the signature.
	token, err := jwt.ParseWithClaims(tokenString, &domain.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(m.conf.Login.Secret), nil
	})

	if err != nil {
		return -1, fmt.Errorf("failed to parse token: %w", err)
	}

	// Verify the token signature, expiration, and permission.
	if !token.Valid {
		return -1, fmt.Errorf("invalid token signature")
	}

	claims, ok := token.Claims.(*domain.CustomClaims)
	if !ok {
		return -1, fmt.Errorf("invalid token claims")
	}

	if claims.ExpiresAt < time.Now().UTC().Unix() {
		return -1, fmt.Errorf("token has expired")
	}

	if !claims.HasPermission(permission) {
		return -1, fmt.Errorf("user does not have permission")
	}

	return claims.UserID, nil
}

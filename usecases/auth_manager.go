package usecases

import (
	"errors"
	"fmt"

	"github.com/bipuldutta/blogzilla/config"
	"github.com/bipuldutta/blogzilla/domain"
	"github.com/golang-jwt/jwt/v5"
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
		return -1, fmt.Errorf("invalid token")
	} else if errors.Is(err, jwt.ErrTokenMalformed) {
		return -1, fmt.Errorf("malformed token")
	} else if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
		// Invalid signature
		return -1, fmt.Errorf("invalid signature")
	} else if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
		// Token is either expired or not active yet
		return -1, fmt.Errorf("expired or inactive token")
	}

	claims, ok := token.Claims.(*domain.CustomClaims)
	if !ok {
		return -1, fmt.Errorf("invalid token claims")
	}

	if !claims.HasPermission(permission) {
		return -1, fmt.Errorf("user does not have permission")
	}

	return claims.UserID, nil
}

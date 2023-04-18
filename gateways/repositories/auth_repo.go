package repositories

import (
	"context"
	"time"

	"github.com/bipuldutta/blogzilla/config"
	"github.com/bipuldutta/blogzilla/domain"
	"github.com/bipuldutta/blogzilla/utils"

	"github.com/golang-jwt/jwt/v5"
)

var authLogger = *utils.Logger()

type AuthRepo struct {
	conf *config.Config
}

func NewAuthRepo(conf *config.Config) domain.AuthRepo {
	return &AuthRepo{
		conf: conf,
	}
}

func (r *AuthRepo) GetToken(ctx context.Context, userID int64, permissions map[string]any) (string, error) {
	currentTime := time.Now().UTC()
	claims := domain.CustomClaims{
		UserID:      userID,
		Permissions: permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			// A usual scenario is to set the expiration time relative to the current time
			ExpiresAt: jwt.NewNumericDate(currentTime.Add(time.Minute * time.Duration(r.conf.Login.Expiry))),
			IssuedAt:  jwt.NewNumericDate(currentTime),
			NotBefore: jwt.NewNumericDate(currentTime),
			Issuer:    "blogzilla",
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := jwtToken.SignedString([]byte(r.conf.Login.Secret))
	if err != nil {
		authLogger.WithError(err).Error("failed to create JWT token")
		return "", err
	}

	return token, nil
}

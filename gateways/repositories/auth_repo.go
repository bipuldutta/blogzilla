package repositories

import (
	"context"
	"time"

	"github.com/bipuldutta/blogzilla/config"
	"github.com/bipuldutta/blogzilla/domain"
	"github.com/bipuldutta/blogzilla/utils"

	"github.com/dgrijalva/jwt-go"
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

func (r *AuthRepo) GetToken(ctx context.Context, userID int64, permissions []string) (string, error) {
	currentTime := time.Now().UTC()
	claims := domain.CustomClaims{
		UserID:      userID,
		Permissions: permissions,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: currentTime.Add(time.Minute * time.Duration(r.conf.Login.Expiry)).Unix(),
			IssuedAt:  currentTime.Unix(),
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

package repositories

import (
	"blogs/config"
	"blogs/domain"
	"blogs/utils"
	"context"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var sessLogger = *utils.Logger()

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
		sessLogger.WithError(err).Error("failed to create JWT token")
		return "", err
	}

	return token, nil
}

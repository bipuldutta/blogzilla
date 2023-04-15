package usecases

import (
	"context"
	"fmt"

	"github.com/bipuldutta/blogzilla/domain"
)

/*
UserManager is the actual business logic section for managing all user related transactions
while this is a skeleton and just making calls to the repo layer at this time there could be
actual BL we could implement at some point
*/
type UserManager struct {
	userRepo domain.UserRepo
}

func NewUserManager(userRepo domain.UserRepo) *UserManager {
	return &UserManager{userRepo: userRepo}
}

func (m *UserManager) Create(ctx context.Context, newUser *domain.User) (*domain.User, error) {
	// validate user input
	if newUser.Username == "" || newUser.Password == "" || newUser.FirstName == "" || newUser.LastName == "" {
		return nil, fmt.Errorf("incomplete user information")
	}

	return m.userRepo.Create(ctx, newUser)
}

func (m *UserManager) Login(ctx context.Context, username string, password string) (string, error) {
	return m.userRepo.Login(ctx, username, password)
}

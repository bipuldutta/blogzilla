package repositories

import (
	"blogs/config"
	"blogs/domain"
	"blogs/utils"
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

const (
	createUserQuery = `INSERT INTO users (username, password, first_name, last_name) VALUES ($1, $2, $3, $4) RETURNING id`
	loginQuery      = `SELECT id, password FROM users WHERE username = $1`
	permissionQuery = `SELECT DISTINCT UNNEST(r.permissions)
		FROM roles r
		JOIN user_roles ur ON r.id = ur.role_id
		WHERE ur.user_id = $1`
)

var userLogger = *utils.Logger()

type UserRepo struct {
	conf        *config.Config
	client      *pgxpool.Pool
	sessionRepo domain.AuthRepo
}

func NewUserRepo(conf *config.Config, client *pgxpool.Pool, sessionRepo domain.AuthRepo) domain.UserRepo {
	return &UserRepo{
		conf:        conf,
		client:      client,
		sessionRepo: sessionRepo,
	}
}

func (r *UserRepo) Create(ctx context.Context, newUser *domain.User) (int64, error) {
	// hash password
	hashedPassword, err := r.hashPassword(newUser.Password)
	if err != nil {
		userLogger.WithError(err).Error("failed to hash password")
		return -1, err
	}
	userLogger.Infof("HASH: %s", hashedPassword)
	var id int64
	// create user and return its id
	err = r.client.QueryRow(ctx, createUserQuery, newUser.Username, hashedPassword, newUser.FirstName, newUser.LastName).Scan(&id)
	if err != nil {
		userLogger.WithError(err).Error("failed to create user")
		return -1, err
	}

	return id, nil
}

func (r *UserRepo) Login(ctx context.Context, username string, password string) (string, error) {

	userLogger.Infof("pwd: %s, HASH: %s", password, password)
	// Get the user from the database
	var userID int64
	var hashedPassword string
	err := r.client.QueryRow(ctx, loginQuery, username).Scan(&userID, &hashedPassword)
	if err != nil {
		return "", fmt.Errorf("invalid username or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return "", fmt.Errorf("invalid username or password")
	}

	permissions, err := r.getUserPermissions(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("failed to get user permissions")
	}
	userLogger.Infof("######: %+v", permissions)
	token, err := r.sessionRepo.GetToken(ctx, userID, permissions)
	if err != nil {
		return "", fmt.Errorf("failed to set the session information")
	}
	return token, nil
}

func (r *UserRepo) getUserPermissions(ctx context.Context, userID int64) ([]string, error) {
	var permissions []string

	rows, err := r.client.Query(ctx, permissionQuery, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var permission string
		if err := rows.Scan(&permission); err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return permissions, nil
}

func (r *UserRepo) hashPassword(password string) (string, error) {
	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error hashing password")
	}
	return string(hashedPassword), nil
}

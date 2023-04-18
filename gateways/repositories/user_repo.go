package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/bipuldutta/blogzilla/config"
	"github.com/bipuldutta/blogzilla/domain"
	"github.com/bipuldutta/blogzilla/utils"

	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

const (
	createUserQuery    = `INSERT INTO users (username, password, first_name, last_name) VALUES ($1, $2, $3, $4) RETURNING id`
	assignUserRoles    = `INSERT INTO user_roles (user_id, role_id) VALUES ($1, $2)`
	getUserByNameQuery = `SELECT id, username, password, first_name, last_name, created_at, updated_at FROM users WHERE username = $1`
	getRoleByName      = `SELECT id, name, description, UNNEST(permissions) FROM roles WHERE name = $1`
	permissionQuery    = `SELECT DISTINCT UNNEST(r.permissions)
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

func (r *UserRepo) Create(ctx context.Context, newUser *domain.User) (*domain.User, error) {
	// hash password
	hashedPassword, err := r.hashPassword(newUser.Password)
	if err != nil {
		userLogger.WithError(err).Error("failed to hash password")
		return nil, err
	}
	var userID int64
	// create user and return its id
	err = r.client.QueryRow(ctx, createUserQuery, newUser.Username, hashedPassword, newUser.FirstName, newUser.LastName).Scan(&userID)
	if err != nil {
		userLogger.WithError(err).Error("failed to create user")
		return nil, err
	}

	// assign editor and viewer roles to this user
	// at this point we are not going to deal with cases like what happens if this transaction fails
	// but in production we should perhaps run these two with a transaction so that we can rollback.
	// Also, we will do a batch query to get both roles in production
	editorRole, err := r.GetRoleByName(ctx, utils.EditorRole)
	if err != nil {
		userLogger.WithError(err).Error("failed to get editor role")
		return nil, err
	}
	viewerRole, err := r.GetRoleByName(ctx, utils.ViewerRole)
	if err != nil {
		userLogger.WithError(err).Error("failed to get viewer role")
		return nil, err
	}
	err = r.AssignRoles(ctx, userID, editorRole.ID, viewerRole.ID)
	if err != nil {
		userLogger.WithError(err).Error("failed to assign roles to the user")
		return nil, err
	}
	// since the create was successful, set it ID
	createdUser, err := r.GetUserByUsername(ctx, newUser.Username)
	if err != nil {
		userLogger.WithError(err).Error("failed to get the user during creation")
		return nil, err
	}

	return createdUser, nil
}

func (r *UserRepo) AssignRoles(ctx context.Context, userID int64, roleIDs ...int64) error {
	// create user and return its id
	for _, roleID := range roleIDs {
		_, err := r.client.Query(ctx, assignUserRoles, userID, roleID)
		if err != nil {
			userLogger.WithError(err).Error("failed to assign roles to user")
			return err
		}
	}
	return nil
}
func (r *UserRepo) GetRoleByName(ctx context.Context, roleName string) (*domain.Role, error) {
	var id int64
	var name, description string
	var permissions []string
	rows, err := r.client.Query(ctx, getRoleByName, roleName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var permission string
		if err := rows.Scan(&id, &name, &description, &permission); err != nil {
			userLogger.WithError(err).Error("failed to get role by name")
			return nil, err
		}
		permissions = append(permissions, permission)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	role := domain.Role{
		ID:          id,
		Name:        name,
		Description: description,
		Permissions: permissions,
	}
	return &role, nil
}

func (r *UserRepo) GetUserByUsername(ctx context.Context, uname string) (*domain.User, error) {
	var userID int64
	var username, password string
	var firstName, lastName sql.NullString
	var createdAt, updatedAt time.Time

	rows, err := r.client.Query(ctx, getUserByNameQuery, uname)
	if err != nil {
		userLogger.WithError(err).Error("failed to check username")
		return nil, fmt.Errorf("failed to check username")
	}
	defer rows.Close()
	if !rows.Next() {
		// user not found
		return nil, nil
	}
	err = rows.Scan(&userID, &username, &password, &firstName, &lastName, &createdAt, &updatedAt)
	if err != nil {
		userLogger.WithError(err).Error("failed to read user data")
		return nil, fmt.Errorf("failed to read user data")
	}

	user := &domain.User{
		ID:        userID,
		Username:  username,
		Password:  password,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
	if firstName.Valid {
		user.FirstName = firstName.String
	}
	if lastName.Valid {
		user.LastName = lastName.String
	}
	return user, nil
}

func (r *UserRepo) Login(ctx context.Context, username string, password string) (string, error) {
	// Get the user from the database
	user, err := r.GetUserByUsername(ctx, username)
	if err != nil {
		return "", fmt.Errorf("failed to get user by name")
	}
	if user == nil {
		return "", fmt.Errorf("user does not exists")
	}

	userID := user.ID
	hashedPassword := user.Password

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return "", fmt.Errorf("invalid username or password")
	}

	permissions, err := r.getUserPermissions(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("failed to get user permissions")
	}
	token, err := r.sessionRepo.GetToken(ctx, userID, permissions)
	if err != nil {
		return "", fmt.Errorf("failed to set the session information")
	}
	return token, nil
}

func (r *UserRepo) getUserPermissions(ctx context.Context, userID int64) (map[string]any, error) {
	permissions := make(map[string]any)

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
		permissions[permission] = nil
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

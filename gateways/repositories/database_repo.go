package repositories

import (
	"context"

	"github.com/bipuldutta/blogzilla/config"
	"github.com/bipuldutta/blogzilla/domain"
	"github.com/bipuldutta/blogzilla/utils"

	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	usersTable = `CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(255) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL,
		first_name VARCHAR(255),
		last_name VARCHAR(255),
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP NOT NULL DEFAULT NOW()
	);`

	blogsTable = `CREATE TABLE IF NOT EXISTS blogs (
	  id SERIAL PRIMARY KEY,
	  user_id INTEGER NOT NULL REFERENCES users(id),
	  title TEXT NOT NULL,
	  content TEXT NOT NULL,
	  tags TEXT,
	  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	  updated_at TIMESTAMP NOT NULL DEFAULT NOW()
	);`

	rolesTable = `CREATE TABLE IF NOT EXISTS roles (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL UNIQUE,
		description TEXT,
		permissions TEXT[]
	);`

	userRolesTable = `CREATE TABLE IF NOT EXISTS user_roles (
	  user_id INT REFERENCES users(id),
	  role_id INT REFERENCES roles(id),
	  PRIMARY KEY (user_id, role_id)
	);`

	rolesData = `INSERT INTO roles (name, description, permissions) VALUES
	('admin', 'Administrator', ARRAY['create_user', 'read_user', 'update_user', 'delete_user', 'create_blog', 'read_blog', 'update_blog', 'delete_blog']),
	('editor', 'Editor', ARRAY['create_blog', 'read_blog', 'update_blog', 'delete_blog']),
	('viewer', 'Viewer', ARRAY['read_user', 'read_blog']) ON CONFLICT DO NOTHING;`
)

var (
	tables = []string{
		usersTable,
		blogsTable,
		rolesTable,
		userRolesTable,
	}
)

var dbLogger = *utils.Logger()

// DatabaseRepo functionalities of this repo is very specific to the initialization of the database tables
// creating roles and an admin user.
type DatabaseRepo struct {
	conf     *config.Config
	client   *pgxpool.Pool
	userRepo domain.UserRepo
}

func NewDatabaseRepo(conf *config.Config, client *pgxpool.Pool, userRepo domain.UserRepo) domain.DatabaseRepo {
	return &DatabaseRepo{
		conf:     conf,
		client:   client,
		userRepo: userRepo,
	}
}

func (r *DatabaseRepo) Initialize(ctx context.Context) error {
	// check and initialize the database tables
	for _, query := range tables {
		dbLogger.Infof("attempting table creation if necessary, query: %s", query)
		_, err := r.client.Exec(ctx, query)
		if err != nil {
			dbLogger.WithError(err).Error("failed to execute query: %s", query)
			return err
		}
	}

	// if necessary create roles with permissions
	_, err := r.client.Exec(ctx, rolesData)
	if err != nil {
		dbLogger.WithError(err).Error("failed to create user")
		return err
	}

	// get the admin role ID
	adminRole, err := r.userRepo.GetRoleByName(ctx, utils.AdminRole)
	if err != nil {
		dbLogger.WithError(err).Error("failed to get admin role")
		return err
	}

	// check if we have already created the admin user, happens during restart
	user, err := r.userRepo.GetUserByUsername(ctx, r.conf.DefaultUser.Username)
	if err != nil {
		dbLogger.WithError(err).Error("failed to check if the default user exists or not")
		return err
	}
	// if the default user does not exists then only attempt creating
	if user == nil {
		adminUser := domain.User{
			Username: r.conf.DefaultUser.Username,
			Password: r.conf.DefaultUser.Password,
		}
		createdAdminUser, err := r.userRepo.Create(ctx, &adminUser)
		if err != nil {
			dbLogger.WithError(err).Error("failed to create admin user")
			return err
		}

		// now assign the admin role id to the admin user
		err = r.userRepo.AssignRoles(ctx, createdAdminUser.ID, adminRole.ID)
		if err != nil {
			dbLogger.WithError(err).Error("failed to assign admin role to the admin user")
			return err
		}
	}
	return nil
}

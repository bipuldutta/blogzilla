package main

import (
	"context"
	"fmt"

	"github.com/bipuldutta/blogzilla/api"
	"github.com/bipuldutta/blogzilla/config"
	"github.com/bipuldutta/blogzilla/gateways/repositories"
	"github.com/bipuldutta/blogzilla/usecases"
	"github.com/bipuldutta/blogzilla/utils"

	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
)

// Create a new instance of the logger. You can have any number of instances.
var logger = *utils.Logger()

/*
This is the main file where our web service will start
*/
func main() {
	ctx := context.Background()

	// Read the config
	conf := config.NewConfig()
	dbPool, err := initDB(ctx, conf.Postgres)
	if err != nil {
		logger.Fatal(err)
	}
	sessionRepo := repositories.NewAuthRepo(conf)
	userRepo := repositories.NewUserRepo(conf, dbPool, sessionRepo)
	userManager := usecases.NewUserManager(userRepo)
	databaseRepo := repositories.NewDatabaseRepo(conf, dbPool, userRepo)
	databaseManager := usecases.NewDatabaseManager(databaseRepo)

	// attempt initializing database tables and default roles, users etc.
	err = databaseManager.Initialize(ctx)
	if err != nil {
		logger.WithError(err).Fatalf("failed to initialize database tables, roles, default user etc.")
	}

	webService := api.NewWebService(conf, userManager)
	err = webService.Start()
	if err != nil {
		logger.WithError(err).Fatalf("failed to start server")
	}

	fmt.Println(conf.Postgres.Database)
}

func initDB(ctx context.Context, config config.PostgresConfig) (*pgxpool.Pool, error) {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", config.User, config.Password, config.Host, config.Port, config.Database)
	dbpool, err := pgxpool.Connect(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}

	logger.Println("successfully connected to database!")
	return dbpool, nil
}

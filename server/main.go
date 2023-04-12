package main

import (
	"blogs/api"
	"blogs/config"
	"blogs/gateways/repositories"
	"blogs/usecases"
	"blogs/utils"
	"context"
	"fmt"

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
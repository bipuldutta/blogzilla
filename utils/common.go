package utils

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

const (
	traceID = "TRACE_ID"

	CreateUserPermission = "create_user"
	ReadUserPermission   = "read_user"
	UpdateUserPermission = "update_user"
	DeleteUserPermission = "delete_user"

	CreateBlogPermission = "create_blog"
	ReadBlogPermission   = "read_blog"
	UpdateBlogPermission = "update_blog"
	DeleteBlogPermission = "delete_blog"
)

// Create a new instance of the logger. You can have any number of instances.
var log *logrus.Logger

func init() {
	log = logrus.New()
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&logrus.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(logrus.InfoLevel)
}

func Logger() *logrus.Logger {
	return log
}

func CreateContext() context.Context {
	ctx := context.Background()
	// add a trace id so that we can put it in the log and see the entire flow
	ctx = context.WithValue(ctx, traceID, getNewGUID())
	return ctx
}

func getNewGUID() string {
	uuidWithHyphen := uuid.New()
	fmt.Println(uuidWithHyphen)
	uuid := strings.Replace(uuidWithHyphen.String(), "-", "", -1)
	return uuid
}

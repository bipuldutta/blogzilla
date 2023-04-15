package utils

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

const (
	traceID = "TRACE_ID"

	AdminRole  = "admin"
	EditorRole = "editor"
	ViewerRole = "viewer"

	CreateUserPermission = "create_user"
	ReadUserPermission   = "read_user"
	UpdateUserPermission = "update_user"
	DeleteUserPermission = "delete_user"

	CreateBlogPermission = "create_blog"
	ReadBlogPermission   = "read_blog"
	UpdateBlogPermission = "update_blog"
	DeleteBlogPermission = "delete_blog"
)

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

package usecases

import (
	"context"

	"github.com/bipuldutta/blogzilla/domain"
)

/*
DatabaseManager is the business logic to initialize database tables, roles and the default user
*/
type DatabaseManager struct {
	databaseRepo domain.DatabaseRepo
}

func NewDatabaseManager(databaseRepo domain.DatabaseRepo) *DatabaseManager {
	return &DatabaseManager{databaseRepo: databaseRepo}
}

func (m *DatabaseManager) Initialize(ctx context.Context) error {
	return m.databaseRepo.Initialize(ctx)
}

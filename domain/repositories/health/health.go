package health

import (
	"context"
	"database/sql"

	"gorm.io/gorm"
)

type (
	IHealthRepositories interface {
		DatabaseHealth(ctx context.Context) (sql.DBStats, error)
	}
	HealthRepositories struct {
		db *gorm.DB
	}
)

func NewHealthRepositories(db *gorm.DB) IHealthRepositories {
	return &HealthRepositories{
		db: db,
	}
}

func (repo HealthRepositories) DatabaseHealth(ctx context.Context) (sql.DBStats, error) {
	sqlDB, _ := repo.db.DB()

	return sqlDB.Stats(), sqlDB.Ping()
}

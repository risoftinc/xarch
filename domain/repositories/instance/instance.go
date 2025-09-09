package instance

import (
	"context"
	"fmt"

	"github.com/risoftinc/xarch/constant"
	"gorm.io/gorm"
)

type (
	IInstanceRepository interface {
		BeginTransaction(ctx context.Context) (*gorm.DB, error)
		BeginTransactionWithContext(ctx context.Context) (context.Context, *gorm.DB, error)
		GetTransactionFromContext(ctx context.Context) (*gorm.DB, bool)
	}

	InstanceRepository struct {
		db *gorm.DB
	}
)

func NewInstanceRepository(db *gorm.DB) IInstanceRepository {
	return &InstanceRepository{
		db: db,
	}
}

// BeginTransaction starts a new database transaction
func (r InstanceRepository) BeginTransaction(ctx context.Context) (*gorm.DB, error) {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}
	return tx, nil
}

// BeginTransactionWithContext starts a new database transaction and puts it in context
func (r InstanceRepository) BeginTransactionWithContext(ctx context.Context) (context.Context, *gorm.DB, error) {
	tx, err := r.BeginTransaction(ctx)
	if err != nil {
		return ctx, nil, err
	}

	return context.WithValue(ctx, constant.TransactionKey, tx), tx, nil
}

// GetTransactionFromContext retrieves transaction from context
func (r InstanceRepository) GetTransactionFromContext(ctx context.Context) (*gorm.DB, bool) {
	tx, ok := ctx.Value(constant.TransactionKey).(*gorm.DB)
	return tx, ok
}

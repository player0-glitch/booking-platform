package database

import (
	"context"

	"gorm.io/gorm"
)

type contextKey string

var txKey contextKey = "gorm_tx"

// godoc
// Interface that defines running attomic database operations
type TransactorInterface interface {
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type gormTransactor struct {
	db *gorm.DB
}

func NewGormTransactor(db *gorm.DB) TransactorInterface {
	return gormTransactor{db: db}
}

// godoc
// Starts a GORM database transaction and inject gorm.DB into it
func (t gormTransactor) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {

	return t.db.Transaction(func(tx *gorm.DB) error {
		//Store the *gorm.DB transaction into the context
		txCtx := context.WithValue(ctx, txKey, tx)
		return fn(txCtx)
	})
}

// godoc
// Get the correct *gorm.DB (are we in a transaction ctx or not)
func DbFromTransaction(ctx context.Context, defaultDb *gorm.DB) *gorm.DB {
	tx, ok := ctx.Value(txKey).(*gorm.DB)
	if ok {
		//If transation is found, run db within it
		return tx
	}
	return defaultDb // fallback to the default db we were already using
}

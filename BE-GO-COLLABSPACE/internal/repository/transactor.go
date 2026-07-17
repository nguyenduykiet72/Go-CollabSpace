package repository

import (
	"context"

	"gorm.io/gorm"
)

type txKey struct{}

type transactor struct {
	db *gorm.DB
}

func NewTransactor(db *gorm.DB) *transactor {
	return &transactor{db}
}

func (t *transactor) WithinTransaction(ctx context.Context, tFunc func(ctx context.Context) error) error {
	return t.db.Transaction(func(tx *gorm.DB) error {
		txCtx := context.WithValue(ctx, txKey{}, tx)
		return tFunc(txCtx)
	})
}

func (t *transactor) getDB(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(txKey{}).(*gorm.DB); ok {
		return tx
	}
	return t.db
}

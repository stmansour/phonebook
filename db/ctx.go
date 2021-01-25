package db

import (
	"context"
	"database/sql"
)

type ctxKey int64

const (
	sessionCtxKey ctxKey = 0
	dbTxCtxKey    ctxKey = iota
)

// SetDBTxContextKey set the session in the given context object
// and returns new context with sql.Tx
func SetDBTxContextKey(ctx context.Context, t *sql.Tx) context.Context {
	return context.WithValue(ctx, dbTxCtxKey, t)
}

// TxFromContext extracts sql.Tx from the given context
// with flag indicating whether sql.Tx found or not
func TxFromContext(ctx context.Context) (*sql.Tx, bool) {
	tx, ok := ctx.Value(dbTxCtxKey).(*sql.Tx)
	return tx, ok
}

// NewTransactionWithContext returns newly created sql.Tx object
// and it also embeds that in the provided context and returns newly updated ctx
func NewTransactionWithContext(ctx context.Context) (*sql.Tx, context.Context, error) {
	var (
		tx  *sql.Tx
		err error
	)

	// get the new transaction
	tx, err = DB.DirDB.Begin()
	if err != nil {
		return tx, ctx, err
	}

	// set the transaction in context
	ctx = SetDBTxContextKey(ctx, tx)

	return tx, ctx, err
}

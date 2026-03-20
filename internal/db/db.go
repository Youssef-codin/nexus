package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TxManager interface {
	StartTx(ctx context.Context) (context.Context, pgx.Tx, error)
}

type DB struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *DB {
	return &DB{pool: pool}
}

func (db *DB) StartTx(ctx context.Context) (context.Context, pgx.Tx, error) {
	tx, err := db.pool.BeginTx(ctx, pgx.TxOptions{})
	txCtx := NewTxContext(ctx, tx)
	return txCtx, tx, err
}

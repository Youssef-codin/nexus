package db

import (
	"context"

	repo "github.com/Youssef-codin/NexusPay/internal/db/postgresql/sqlc"
	"github.com/jackc/pgx/v5"
)

type ctxKeyTx struct{}

func Queries(ctx context.Context, pool *pgx.Conn) *repo.Queries {
	if tx, ok := ctx.Value(ctxKeyTx{}).(pgx.Tx); ok {
		return repo.New(tx)
	}
	return repo.New(pool)
}

func NewTxContext(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, ctxKeyTx{}, tx)
}

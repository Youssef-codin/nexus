package transactions

import (
	"context"

	"github.com/Youssef-codin/NexusPay/internal/db"
	repo "github.com/Youssef-codin/NexusPay/internal/db/postgresql/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

type transactionRepo interface {
	CreateTransaction(
		ctx context.Context,
		arg repo.CreateTransactionParams,
	) (repo.CreateTransactionRow, error)
	GetTransactionById(ctx context.Context, id pgtype.UUID) (repo.Transaction, error)
	GetTransactionByTransferId(
		ctx context.Context,
		transferID pgtype.UUID,
	) (repo.Transaction, error)
	GetTransactionsByWalletId(ctx context.Context, walletID pgtype.UUID) ([]repo.Transaction, error)
	UpdateTransactionStatus(
		ctx context.Context,
		arg repo.UpdateTransactionStatusParams,
	) (repo.Transaction, error)
}

type TransactionRepo struct {
	db *db.DB
}

func NewTransactionRepo(database *db.DB) transactionRepo {
	return &TransactionRepo{db: database}
}

func (r *TransactionRepo) CreateTransaction(
	ctx context.Context,
	arg repo.CreateTransactionParams,
) (repo.CreateTransactionRow, error) {
	return db.Queries(ctx, r.db).CreateTransaction(ctx, arg)
}

func (r *TransactionRepo) GetTransactionById(
	ctx context.Context,
	id pgtype.UUID,
) (repo.Transaction, error) {
	return db.Queries(ctx, r.db).GetTransactionById(ctx, id)
}

func (r *TransactionRepo) GetTransactionByTransferId(
	ctx context.Context,
	transferID pgtype.UUID,
) (repo.Transaction, error) {
	return db.Queries(ctx, r.db).GetTransactionByTransferId(ctx, transferID)
}

func (r *TransactionRepo) GetTransactionsByWalletId(
	ctx context.Context,
	walletID pgtype.UUID,
) ([]repo.Transaction, error) {
	return db.Queries(ctx, r.db).GetTransactionsByWalletId(ctx, walletID)
}

func (r *TransactionRepo) UpdateTransactionStatus(
	ctx context.Context,
	arg repo.UpdateTransactionStatusParams,
) (repo.Transaction, error) {
	return db.Queries(ctx, r.db).UpdateTransactionStatus(ctx, arg)
}

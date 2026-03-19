package stripe

import (
	"context"
	"errors"

	dbpkg "github.com/Youssef-codin/NexusPay/internal/db"
	repo "github.com/Youssef-codin/NexusPay/internal/db/postgresql/sqlc"
	"github.com/Youssef-codin/NexusPay/internal/transactions"
	"github.com/Youssef-codin/NexusPay/internal/wallet"
	"github.com/jackc/pgx/v5"
)

var (
	ErrParseEvent        = errors.New("failed to parse event")
	ErrAlreadyProcessing = errors.New("already processing transaction")
)

type IService interface {
	HandlePaymentSucceeded(ctx context.Context, req HandlePaymentSucceededRequest) error
	HandlePaymentFailed(ctx context.Context, req HandlePaymentFailedRequest) error
	HandlePaymentCanceled(ctx context.Context, req HandlePaymentCanceledRequest) error
}

type WebhookService struct {
	pool           *pgx.Conn
	walletSvc      wallet.IService
	transactionSvc transactions.IService
}

func NewWebhookService(
	pool *pgx.Conn,
	walletSvc wallet.IService,
	transactionSvc transactions.IService,
) IService {
	return &WebhookService{
		pool:           pool,
		walletSvc:      walletSvc,
		transactionSvc: transactionSvc,
	}
}

func (svc *WebhookService) HandlePaymentSucceeded(
	ctx context.Context,
	req HandlePaymentSucceededRequest,
) error {
	tx, err := svc.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	txCtx := dbpkg.NewTxContext(ctx, tx)

	transaction, err := svc.transactionSvc.GetById(txCtx, transactions.GetByIdRequest{
		ID: req.TransactionID,
	})

	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			svc.transactionSvc.UpdateStatus(ctx, transactions.UpdateTransactionRequest{
				ID:     transaction.ID,
				Status: repo.TransactionStatusFailed,
			})
		}
	}()

	if transaction.Status == repo.TransactionStatusProcessing {
		return ErrAlreadyProcessing
	}

	err = svc.transactionSvc.UpdateStatus(txCtx, transactions.UpdateTransactionRequest{
		ID:     transaction.ID,
		Status: repo.TransactionStatusProcessing,
	})

	if err != nil {
		return err
	}

	_, err = svc.walletSvc.AddToWallet(txCtx, wallet.AddToWalletRequest{
		WalletID: transaction.WalletID,
		Amount:   transaction.Amount,
	})

	if err != nil {
		return err
	}

	err = svc.transactionSvc.UpdateStatus(
		txCtx,
		transactions.UpdateTransactionRequest{
			ID:     transaction.ID,
			Status: repo.TransactionStatusCompleted,
		},
	)

	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (svc *WebhookService) HandlePaymentFailed(
	ctx context.Context,
	req HandlePaymentFailedRequest,
) error {
	return svc.transactionSvc.UpdateStatus(ctx, transactions.UpdateTransactionRequest{
		ID:     req.TransactionID,
		Status: repo.TransactionStatusFailed,
	})
}

func (svc *WebhookService) HandlePaymentCanceled(
	ctx context.Context,
	req HandlePaymentCanceledRequest,
) error {
	return svc.transactionSvc.UpdateStatus(ctx, transactions.UpdateTransactionRequest{
		ID:     req.TransactionID,
		Status: repo.TransactionStatusFailed,
	})
}

package stripe

import (
	"context"
	"errors"
	"testing"

	"github.com/Youssef-codin/NexusPay/internal/transactions"
	"github.com/Youssef-codin/NexusPay/internal/wallet"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockWalletSvc struct {
	mock.Mock
}

func (m *MockWalletSvc) GetById(ctx context.Context, req wallet.GetWalletRequest) (wallet.GetWalletResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(wallet.GetWalletResponse), args.Error(1)
}

func (m *MockWalletSvc) GetByUserId(ctx context.Context) (wallet.GetWalletResponse, error) {
	args := m.Called(ctx)
	return args.Get(0).(wallet.GetWalletResponse), args.Error(1)
}

func (m *MockWalletSvc) CreateWallet(ctx context.Context, req wallet.CreateWalletRequest) (wallet.CreateWalletResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(wallet.CreateWalletResponse), args.Error(1)
}

func (m *MockWalletSvc) TopUp(ctx context.Context, req wallet.TopUpRequest) (wallet.TopUpResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(wallet.TopUpResponse), args.Error(1)
}

func (m *MockWalletSvc) DeductFromBalance(ctx context.Context, req wallet.DeductRequest) (wallet.DeductResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(wallet.DeductResponse), args.Error(1)
}

func (m *MockWalletSvc) AddToWallet(ctx context.Context, req wallet.AddToWalletRequest) (wallet.AddToWalletResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(wallet.AddToWalletResponse), args.Error(1)
}

type MockTransactionSvc struct {
	mock.Mock
}

func (m *MockTransactionSvc) GetById(ctx context.Context, req transactions.GetByIdRequest) (transactions.GetTransactionResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(transactions.GetTransactionResponse), args.Error(1)
}

func (m *MockTransactionSvc) CreateTransaction(ctx context.Context, req transactions.CreateTransactionRequest) (transactions.CreateTransactionResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(transactions.CreateTransactionResponse), args.Error(1)
}

func (m *MockTransactionSvc) UpdateStatus(ctx context.Context, req transactions.UpdateTransactionRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func TestHandlePaymentSucceeded_Integration(t *testing.T) {
	t.Skip("Integration test - requires database connection and real pool.BeginTx")
}

func TestHandlePaymentFailed(t *testing.T) {
	ctx := context.Background()
	transactionID := uuid.New()

	t.Run("success", func(t *testing.T) {
		mockWalletSvc := new(MockWalletSvc)
		mockTxSvc := new(MockTransactionSvc)

		mockTxSvc.On("UpdateStatus", mock.Anything, mock.Anything).Return(nil)

		svc := &WebhookService{
			walletSvc:      mockWalletSvc,
			transactionSvc: mockTxSvc,
		}

		err := svc.HandlePaymentFailed(ctx, HandlePaymentFailedRequest{
			TransactionID: transactionID.String(),
		})

		assert.NoError(t, err)
		mockTxSvc.AssertExpectations(t)
	})

	t.Run("transaction_not_found", func(t *testing.T) {
		mockWalletSvc := new(MockWalletSvc)
		mockTxSvc := new(MockTransactionSvc)

		mockTxSvc.On("UpdateStatus", mock.Anything, mock.Anything).Return(transactions.ErrTransactionNotFound)

		svc := &WebhookService{
			walletSvc:      mockWalletSvc,
			transactionSvc: mockTxSvc,
		}

		err := svc.HandlePaymentFailed(ctx, HandlePaymentFailedRequest{
			TransactionID: transactionID.String(),
		})

		assert.ErrorIs(t, err, transactions.ErrTransactionNotFound)
		mockTxSvc.AssertExpectations(t)
	})

	t.Run("database_error", func(t *testing.T) {
		mockWalletSvc := new(MockWalletSvc)
		mockTxSvc := new(MockTransactionSvc)

		mockTxSvc.On("UpdateStatus", mock.Anything, mock.Anything).Return(errors.New("database error"))

		svc := &WebhookService{
			walletSvc:      mockWalletSvc,
			transactionSvc: mockTxSvc,
		}

		err := svc.HandlePaymentFailed(ctx, HandlePaymentFailedRequest{
			TransactionID: transactionID.String(),
		})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")
		mockTxSvc.AssertExpectations(t)
	})
}

func TestHandlePaymentCanceled(t *testing.T) {
	ctx := context.Background()
	transactionID := uuid.New()

	t.Run("success", func(t *testing.T) {
		mockWalletSvc := new(MockWalletSvc)
		mockTxSvc := new(MockTransactionSvc)

		mockTxSvc.On("UpdateStatus", mock.Anything, mock.Anything).Return(nil)

		svc := &WebhookService{
			walletSvc:      mockWalletSvc,
			transactionSvc: mockTxSvc,
		}

		err := svc.HandlePaymentCanceled(ctx, HandlePaymentCanceledRequest{
			TransactionID: transactionID.String(),
		})

		assert.NoError(t, err)
		mockTxSvc.AssertExpectations(t)
	})

	t.Run("transaction_not_found", func(t *testing.T) {
		mockWalletSvc := new(MockWalletSvc)
		mockTxSvc := new(MockTransactionSvc)

		mockTxSvc.On("UpdateStatus", mock.Anything, mock.Anything).Return(transactions.ErrTransactionNotFound)

		svc := &WebhookService{
			walletSvc:      mockWalletSvc,
			transactionSvc: mockTxSvc,
		}

		err := svc.HandlePaymentCanceled(ctx, HandlePaymentCanceledRequest{
			TransactionID: transactionID.String(),
		})

		assert.ErrorIs(t, err, transactions.ErrTransactionNotFound)
		mockTxSvc.AssertExpectations(t)
	})

	t.Run("database_error", func(t *testing.T) {
		mockWalletSvc := new(MockWalletSvc)
		mockTxSvc := new(MockTransactionSvc)

		mockTxSvc.On("UpdateStatus", mock.Anything, mock.Anything).Return(errors.New("database error"))

		svc := &WebhookService{
			walletSvc:      mockWalletSvc,
			transactionSvc: mockTxSvc,
		}

		err := svc.HandlePaymentCanceled(ctx, HandlePaymentCanceledRequest{
			TransactionID: transactionID.String(),
		})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")
		mockTxSvc.AssertExpectations(t)
	})
}

func TestIServiceInterface(t *testing.T) {
	var _ IService = (*WebhookService)(nil)
}

func TestIServiceMethods(t *testing.T) {
	t.Run("HandlePaymentSucceeded", func(t *testing.T) {
		var svc *WebhookService
		var _ func(context.Context, HandlePaymentSucceededRequest) error = svc.HandlePaymentSucceeded
	})

	t.Run("HandlePaymentFailed", func(t *testing.T) {
		var svc *WebhookService
		var _ func(context.Context, HandlePaymentFailedRequest) error = svc.HandlePaymentFailed
	})

	t.Run("HandlePaymentCanceled", func(t *testing.T) {
		var svc *WebhookService
		var _ func(context.Context, HandlePaymentCanceledRequest) error = svc.HandlePaymentCanceled
	})
}

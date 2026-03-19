package stripe

import (
	"context"
	"errors"
	"log/slog"

	"github.com/Youssef-codin/NexusPay/internal/payment"
	"github.com/stripe/stripe-go/v84"
	"github.com/stripe/stripe-go/v84/paymentintent"
	"github.com/stripe/stripe-go/v84/refund"
)

var (
	ErrPaymentFailed    = errors.New("payment failed")
	ErrRefundFailed     = errors.New("refund failed")
	ErrPaymentCancelled = errors.New("error cancelld")
)

type Service struct{}

func NewService(apiKey string) payment.IService {
	stripe.Key = apiKey
	return &Service{}
}

func (svc *Service) ProcessPayment(
	ctx context.Context,
	req payment.ProcessPaymentRequest,
) (payment.ProcessPaymentResponse, error) {

	params := &stripe.PaymentIntentParams{
		Amount:      new(req.Amount),
		Currency:    stripe.String("egp"),
		Description: new(req.Description),
		Metadata: map[string]string{
			"transaction_id": req.TransactionID,
		},
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: new(true),
		},
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		slog.Error("stripe payment failed", "error", err)
		return payment.ProcessPaymentResponse{}, ErrPaymentFailed
	}

	return payment.ProcessPaymentResponse{
		ProviderPaymentID: pi.ID,
		Status:            FromStripePaymentIntentStatus(pi.Status),
		ClientSecret:      pi.ClientSecret,
	}, nil
}

func (svc *Service) Refund(
	ctx context.Context,
	req payment.RefundRequest,
) (payment.RefundResponse, error) {
	params := &stripe.RefundParams{
		PaymentIntent: stripe.String(req.ProviderPaymentID),
	}

	if req.Amount != nil {
		params.Amount = new(*req.Amount)
	}

	r, err := refund.New(params)
	if err != nil {
		slog.Error("stripe payment failed", "error", err)
		return payment.RefundResponse{}, ErrRefundFailed
	}

	return payment.RefundResponse{
		RefundID: r.ID,
		Status:   FromStripeRefundStatus(r.Status),
		Amount:   r.Amount,
	}, nil
}

func (svc *Service) CancelPayment(
	ctx context.Context,
	req payment.CancelPaymentRequest,
) (payment.CancelPaymentResponse, error) {
	pi, err := paymentintent.Cancel(req.ProviderPaymentID, nil)
	if err != nil {
		slog.Error("stripe payment failed", "error", err)
		return payment.CancelPaymentResponse{}, ErrPaymentCancelled
	}

	return payment.CancelPaymentResponse{
		ProviderPaymentID: pi.ID,
		Status:            FromStripePaymentIntentStatus(pi.Status),
	}, nil
}

func FromStripePaymentIntentStatus(s stripe.PaymentIntentStatus) payment.PaymentStatus {
	switch s {
	case stripe.PaymentIntentStatusSucceeded:
		return payment.PaymentStatusCompleted
	case stripe.PaymentIntentStatusCanceled,
		stripe.PaymentIntentStatusRequiresPaymentMethod:
		return payment.PaymentStatusFailed
	case stripe.PaymentIntentStatusProcessing,
		stripe.PaymentIntentStatusRequiresAction,
		stripe.PaymentIntentStatusRequiresCapture,
		stripe.PaymentIntentStatusRequiresConfirmation:
		return payment.PaymentStatusPending
	default:
		return payment.PaymentStatusPending
	}
}

func FromStripeRefundStatus(s stripe.RefundStatus) payment.PaymentStatus {
	switch s {
	case stripe.RefundStatusSucceeded:
		return payment.PaymentStatusCompleted
	case stripe.RefundStatusFailed,
		stripe.RefundStatusCanceled:
		return payment.PaymentStatusFailed
	case stripe.RefundStatusPending,
		stripe.RefundStatusRequiresAction:
		return payment.PaymentStatusPending
	default:
		return payment.PaymentStatusPending
	}
}

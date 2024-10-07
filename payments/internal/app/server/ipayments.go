package grpcserver

import (
	"context"
	"github.com/NStegura/saga/payments/internal/services/payment/models"
)

// Payments интерфейс для работы с бизнес слоем.
type Payments interface {
	CreatePayment(ctx context.Context, orderID int64) (id int64, err error)
	UpdatePaymentStatus(ctx context.Context, orderID int64, status models.PaymentMessageStatus) (err error)
}

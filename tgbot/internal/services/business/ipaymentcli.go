package business

import (
	"context"
)

type IPaymentCli interface {
	PayOrder(ctx context.Context, orderID int64, status bool) (err error)
}

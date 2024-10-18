package consumer

import (
	"context"
)

// Payments интерфейс для работы с бизнес слоем.
type Payments interface {
	CreatePayment(ctx context.Context, orderID int64) (id int64, err error)
}

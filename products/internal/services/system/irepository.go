package system

import (
	"context"
)

type Repository interface {
	Ping(ctx context.Context) error
}

package product

import (
	"context"

	"github.com/google/uuid"
)

type Cache interface {
	Get(context.Context, uuid.UUID) error
	Set(context.Context, uuid.UUID) error
}

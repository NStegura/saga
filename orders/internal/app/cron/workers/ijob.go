package workers

import (
	"context"
	"time"
)

type Worker interface {
	GetFrequency() time.Duration
	Run(ctx context.Context) error
}

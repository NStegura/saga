package jobs

import (
	"context"
	"time"
)

type Job interface {
	GetFrequency() time.Duration
	Run(ctx context.Context) error
}

package event

import (
	"github.com/NStegura/saga/golibs/event/sender"
	"time"
)

type Option func(*sender.PushEventsWorker)

func WithFrequency(f time.Duration) Option {
	return func(w *sender.PushEventsWorker) {
		w.Frequency = f
		return
	}
}

func WithRateLimit(r int) Option {
	return func(w *sender.PushEventsWorker) {
		w.RateLimit = r
		return
	}
}

func WithEventsLimit(e int) Option {
	return func(w *sender.PushEventsWorker) {
		w.EventsLimit = e
		return
	}
}

func WithReserve(r time.Duration) Option {
	return func(w *sender.PushEventsWorker) {
		w.Reserve = r
		return
	}
}

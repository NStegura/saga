package event

import (
	"github.com/NStegura/saga/golibs/event/repo"
	"github.com/NStegura/saga/golibs/event/sender"
	"github.com/NStegura/saga/golibs/event/service"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	defaultFrequency   = 5 * time.Second
	defaultRateLimit   = 1
	defaultEventsLimit = 50
	defaultReserve     = 60 * time.Second
)

// NewPushEventsWorker иницилизирует задачу
// по получению эвентов для дальнейшей отправки.
func NewPushEventsWorker(
	producer sender.Producer,
	event sender.Event,
	logger *logrus.Logger,
	options ...Option,
) *sender.PushEventsWorker {
	s := &sender.PushEventsWorker{
		Frequency:   defaultFrequency,
		RateLimit:   defaultRateLimit,
		EventsLimit: defaultEventsLimit,
		Reserve:     defaultReserve,
		Producer:    producer,
		Event:       event,
		Logger:      logger,
	}
	for _, opt := range options {
		opt(s)
	}
	return s
}

// NewEventService инициализует бизнес слой.
func NewEventService(
	repo service.Repository,
	logger *logrus.Logger,
) *service.Event {
	return &service.Event{
		Repo:   repo,
		Logger: logger,
	}
}

// NewEventRepository инициализует слой по работе с бд.
func NewEventRepository(
	logger *logrus.Logger,
) repo.EventRepo {
	return repo.EventRepo{
		Logger: logger,
	}
}

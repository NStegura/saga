package event

import (
	"events/repo"
	events "events/sender"
	"events/service"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"time"
)

// NewPushEventsWorker иницилизирует задачу
// по получению эвентов для дальнейшей отправки.
func NewPushEventsWorker(
	frequency time.Duration,
	rateLimit int,
	eventsLimit int,
	reserve time.Duration,
	producer events.Producer,
	event events.Event,
	logger *logrus.Logger,
) *events.PushEventsWorker {
	return &events.PushEventsWorker{
		Frequency:   frequency,
		RateLimit:   rateLimit,
		EventsLimit: eventsLimit,
		Reserve:     reserve,
		Producer:    producer,
		Event:       event,
		Logger:      logger,
	}
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
	pool *pgxpool.Pool,
	logger *logrus.Logger,
) *repo.EventRepository {
	return &repo.EventRepository{
		Pool:   pool,
		Logger: logger,
	}
}

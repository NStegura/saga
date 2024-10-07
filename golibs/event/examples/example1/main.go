package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/NStegura/saga/golibs/event"
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

const (
	dbDSN           = "postgres://usr:psswrd@localhost:54321/example?sslmode=disable"
	pushFrequency   = 5 * time.Second
	pushRateLimit   = 1
	pushEventsLimit = 50
	pushReserveTime = time.Second * 60
)

type SenderStdOut struct {
	log *logrus.Logger
}

func (s *SenderStdOut) PushMsg(msg []byte, topic string) error {
	s.log.Infof("[Pusher][topic:%s] msg: %s", topic, msg)
	return nil
}

func main() {
	ctx, cancelCtx := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancelCtx()

	g, ctx := errgroup.WithContext(ctx)

	log := logrus.New()
	pool, err := pgxpool.New(ctx, dbDSN)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		log.Info("close pool")
		pool.Close()
		log.Info("pool closed")
	}()

	query := `
	BEGIN;
	DO $$
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'event_status') THEN
		CREATE TYPE event_status AS ENUM ('WAIT', 'DONE');
		END IF;
	END
	$$;
	CREATE TABLE IF NOT EXISTS event
	(
		id          bigserial PRIMARY KEY,
		topic       TEXT NOT NULL,
		payload     JSONB NOT NULL,
		status      event_status NOT NULL DEFAULT 'WAIT',
		created_at  timestamp NOT NULL DEFAULT NOW(),
		reserved_to timestamp DEFAULT NULL
	);
	CREATE INDEX IF NOT EXISTS idx_created_at ON "event"(created_at);
	COMMIT;
`
	_, err = pool.Exec(ctx, query)
	if err != nil {
		log.Fatal(err)
	}

	repo := event.NewEventRepository(pool, log)
	service := event.NewEventService(repo, log)
	worker := event.NewPushEventsWorker(
		&SenderStdOut{log: log},
		service,
		log,
		event.WithFrequency(pushFrequency),
		event.WithRateLimit(pushRateLimit),
		event.WithEventsLimit(pushEventsLimit),
		event.WithReserve(pushReserveTime),
	)

	// pusher
	g.Go(func() (err error) {
		pushEventsTimer := time.NewTicker(worker.GetFrequency())
		defer pushEventsTimer.Stop()
		i := 0
		for {
			select {
			case <-pushEventsTimer.C:
				i++
				log.Infof("[Worker|%v] events push", i)
				if err := worker.Run(ctx); err != nil {
					log.Error("Worker Run failed %v", err)
					continue
				}
			case <-ctx.Done():
				return
			}
		}
	})

	// db saver
	g.Go(func() (err error) {
		saveToDBEventsTimer := time.NewTicker(worker.GetFrequency() / 10)
		defer saveToDBEventsTimer.Stop()
		i := 0
		for {
			select {
			case <-saveToDBEventsTimer.C:
				i++
				log.Infof("[saveToDBEventsTimer|%v] events save", i)

				tx, err := repo.OpenTransaction(ctx)
				if err != nil {
					return fmt.Errorf("failed to open transaction: %w", err)
				}

				payload, err := json.Marshal(map[string]any{"number": i, "hello": "world"})
				if err != nil {
					return fmt.Errorf("failed to marshal message: %w", err)
				}
				if err := repo.CreateEvent(ctx, tx, "topic_name", payload); err != nil {
					log.Error(err)
					continue
				}
				_ = repo.Commit(ctx, tx)
			case <-ctx.Done():
				return
			}
		}
	})

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}

}

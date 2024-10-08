package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/NStegura/saga/golibs/event/repo/models"
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
	"time"
)

type EventRepo struct {
	logger *logrus.Logger
}

func New(logger *logrus.Logger) EventRepo {
	return EventRepo{logger: logger}
}

func (e *EventRepo) CreateEvent(
	ctx context.Context,
	tx pgx.Tx,
	topic string,
	payload json.RawMessage,
) (err error) {
	var id int64
	const query = `
		INSERT INTO "event" (topic, payload) 
		VALUES ($1, $2) 
		RETURNING id;
	`

	err = tx.QueryRow(ctx, query,
		topic,
		payload,
	).Scan(&id)

	if err != nil {
		return fmt.Errorf("CreateEvent failed, %w", err)
	}
	e.Logger.Debugf("Create event, id, %v", id)
	return
}

func (e *EventRepo) GetNotProcessedEvents(
	ctx context.Context,
	tx pgx.Tx,
	limit int64,
) (messages []models.EventEntry, err error) {
	var rows pgx.Rows

	const query = `
		SELECT id, payload, topic, status, created_at
		FROM "event"
		WHERE (status = 'WAIT' OR reserved_to < NOW())
		ORDER BY created_at ASC 
		LIMIT $1
		FOR UPDATE SKIP LOCKED;
	`

	rows, err = tx.Query(ctx, query, limit)
	if err != nil {
		return messages, fmt.Errorf("failed to get events: %w", err)
	}

	for rows.Next() {
		var o models.EventEntry
		err = rows.Scan(
			&o.ID,
			&o.Payload,
			&o.Topic,
			&o.Status,
			&o.CreatedAt,
		)
		if err != nil {
			return messages, fmt.Errorf("failed to get events: %w", err)
		}
		messages = append(messages, o)
	}

	if err = rows.Err(); err != nil {
		return messages, fmt.Errorf("failed to get events: %w", err)
	}

	return messages, nil
}

func (e *EventRepo) UpdateReservedTimeEvents(
	ctx context.Context,
	tx pgx.Tx,
	eventsIDs []int64,
	reservedTo time.Time,
) (err error) {
	const query = `
		UPDATE event
        SET reserved_to = $2
        WHERE id = ANY($1)
	`
	_, err = tx.Exec(ctx, query, eventsIDs, reservedTo)
	if err != nil {
		return fmt.Errorf("UpdateReservedTimeEvents failed, %w", err)
	}
	return nil
}

func (e *EventRepo) UpdateEventStatusToDone(
	ctx context.Context,
	tx pgx.Tx,
	eventID int64,
) (err error) {
	const query = `
		UPDATE event
        SET status = 'DONE', reserved_to = NULL
        WHERE id = $1
	`
	_, err = tx.Exec(ctx, query, eventID)
	if err != nil {
		return fmt.Errorf("UpdateEventStatusToDone failed, %w", err)
	}
	return nil
}

package repo

import (
	"context"
	"encoding/json"
	"events/repo/models"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"time"
)

type EventRepository struct {
	Pool *pgxpool.Pool

	Logger *logrus.Logger
}

func (e *EventRepository) CreateEvent(
	ctx context.Context,
	tx pgx.Tx,
	eventType string,
	payload json.RawMessage,
) (err error) {
	var id int64
	const query = `
		INSERT INTO "event" (event_type, payload) 
		VALUES ($1, $2) 
		RETURNING id;
	`

	err = tx.QueryRow(ctx, query,
		eventType,
		payload,
	).Scan(&id)

	if err != nil {
		return fmt.Errorf("CreateEvent failed, %w", err)
	}
	e.Logger.Debugf("Create event, id, %v", id)
	return
}

func (e *EventRepository) GetNotProcessedEvents(
	ctx context.Context,
	tx pgx.Tx,
	limit int64,
) (messages []models.EventEntry, err error) {
	var rows pgx.Rows

	const query = `
		SELECT id, payload, topic, status
		FROM "event"
		WHERE (status = 'WAIT' OR reserved_to < NOW())
		ORDER BY created_at ASC 
		LIMIT $1
		FOR UPDATE SKIP LOCKED;
	`

	rows, err = tx.Query(ctx, query, limit)
	if err != nil {
		return messages, fmt.Errorf("get events failed, %w", err)
	}

	for rows.Next() {
		var o models.EventEntry
		err = rows.Scan(
			&o.ID,
			&o.Payload,
			&o.Status,
		)
		if err != nil {
			return messages, fmt.Errorf("get events failed, %w", err)
		}
		messages = append(messages, o)
	}

	if err = rows.Err(); err != nil {
		return messages, fmt.Errorf("get events failed, %w", err)
	}

	return messages, nil
}

func (e *EventRepository) UpdateReservedTimeEvents(
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

func (e *EventRepository) UpdateEventStatusToDone(
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

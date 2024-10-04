package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/NStegura/saga/payments/internal/repo/models"
	"github.com/jackc/pgx/v5"
)

func (db *DB) CreateOutbox(
	ctx context.Context,
	tx pgx.Tx,
	payload json.RawMessage,
) (err error) {
	var id int64
	const query = `
		INSERT INTO "outbox" (payload, status) 
		VALUES ($1, $2) 
		RETURNING id;
	`

	err = tx.QueryRow(ctx, query,
		payload,
		models.WAIT,
	).Scan(&id)

	if err != nil {
		return fmt.Errorf("CreateOutbox failed, %w", err)
	}
	db.logger.Debugf("Create outbox, id, %v", id)
	return
}

func (db *DB) GetNotProcessedEvents(
	ctx context.Context,
	tx pgx.Tx,
	limit int64,
) (messages []models.OutboxEntry, err error) {
	var rows pgx.Rows

	const query = `
		SELECT id, payload, status
		FROM "outbox"
		WHERE status = 'WAIT' 
		ORDER BY created_at ASC 
		LIMIT $1 
		FOR UPDATE SKIP LOCKED;
	`

	rows, err = tx.Query(ctx, query, limit)
	if err != nil {
		return messages, fmt.Errorf("get events failed, %w", err)
	}

	for rows.Next() {
		var o models.OutboxEntry
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

func (db *DB) UpdateOutboxEvents(ctx context.Context, tx pgx.Tx, messageIDs []int64) (err error) {
	const query = `
		UPDATE outbox 
        SET status = 'DONE' 
        WHERE id = ANY($1)
	`
	_, err = tx.Exec(ctx, query, messageIDs)
	if err != nil {
		return fmt.Errorf("UpdateOutboxEvents failed, %w", err)
	}
	return nil
}

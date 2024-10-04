-- +goose Up
-- +goose StatementBegin

BEGIN;
CREATE TYPE outbox_status AS ENUM ('WAIT', 'PROCESSED');
CREATE TABLE IF NOT EXISTS outbox
(
    id          bigserial PRIMARY KEY,
    payload     JSONB NOT NULL,
    status      outbox_status NOT NULL,
    created_at  timestamp NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_created_at ON "outbox"(created_at);
COMMIT;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP EXTENSION IF EXISTS postgres_fdw;

-- +goose StatementEnd

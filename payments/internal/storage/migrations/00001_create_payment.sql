-- +goose Up
-- +goose StatementBegin
BEGIN;
CREATE TYPE payment_status AS ENUM ('CREATED', 'FAILED', 'COMPLETED');
CREATE TABLE IF NOT EXISTS payment
(
    id          bigserial PRIMARY KEY,
	order_id    bigint UNIQUE NOT NULL,
	status      payment_status NOT NULL,
    created_at  timestamp NOT NULL DEFAULT NOW(),
    updated_at  timestamp NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_payment_created_at ON "payment"(created_at);
COMMIT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_created_at;
DROP TABLE IF EXISTS payment;
DROP TYPE IF EXISTS status_type;

-- +goose StatementEnd

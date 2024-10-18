-- +goose Up
-- +goose StatementBegin
BEGIN;
CREATE TYPE order_state AS ENUM (
    'ORDER_CREATED',
    'RESERVE_CREATED',
    'RESERVE_FAILED',
    'PAYMENT_CREATED',
    'PAYMENT_FAILED',
    'PAYMENT_COMPLETED'
);
CREATE TABLE IF NOT EXISTS state
(
    id         bigserial PRIMARY KEY,
    order_id   bigint NOT NULL,
    state      order_state NOT NULL,
    created_at timestamp NOT NULL DEFAULT NOW(),
    CONSTRAINT FK_order_state FOREIGN KEY(order_id) REFERENCES "order"(id)
);
COMMIT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS state;

-- +goose StatementEnd

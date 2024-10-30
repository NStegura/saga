-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS "order"
(
	id               bigserial PRIMARY KEY,
    user_id          bigint NOT NULL,
	description      text
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS "order";

-- +goose StatementEnd

-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS product
(
	id          bigserial PRIMARY KEY,
	category    text NOT NULL,
	name        text NOT NULL,
	description text,
	count       bigint NOT NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS product;

-- +goose StatementEnd

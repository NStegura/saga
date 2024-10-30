-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS product
(
    id         bigserial PRIMARY KEY,
    order_id   bigint NOT NULL,
    product_id bigint NOT NULL,
    count      bigint NOT NULL,
    CONSTRAINT FK_order_product FOREIGN KEY(order_id) REFERENCES "order"(id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS product;

-- +goose StatementEnd

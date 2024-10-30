-- +goose Up
-- +goose StatementBegin

CREATE TABLE reserve
(
	id         bigserial PRIMARY KEY,
	product_id bigint NOT NULL,
	order_id   bigint NOT NULL,
	count      bigint NOT NULL,
	pay_status bool NULL,
    saved_at   timestamp NOT NULL DEFAULT NOW(),
    CONSTRAINT FK_reserve_product FOREIGN KEY(product_id) REFERENCES product(id)
                                                                    ON DELETE RESTRICT
                                                                    ON UPDATE CASCADE
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS reserve;

-- +goose StatementEnd

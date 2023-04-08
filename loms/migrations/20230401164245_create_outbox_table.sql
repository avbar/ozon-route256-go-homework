-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS orders_outbox (
    order_id bigint,
    status varchar,
    created_at timestamp,
    PRIMARY KEY (order_id, status)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS orders_outbox;
-- +goose StatementEnd

-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS orders (
    order_id bigserial PRIMARY KEY,
    status varchar NOT NULL,
    user_id bigint NOT NULL,
    created_at timestamp,
    updated_at timestamp
);

CREATE TABLE IF NOT EXISTS order_items (
    order_id bigint,
    sku integer,
    count smallint,
    PRIMARY KEY (order_id, sku)
);

CREATE TABLE IF NOT EXISTS stocks (
    sku integer,
    warehouse_id bigint,
    count bigint NOT NULL,
    PRIMARY KEY (sku, warehouse_id)
);

CREATE TABLE IF NOT EXISTS reserves (
    sku integer,
    warehouse_id bigint,
    order_id bigint,
    count bigint NOT NULL,
    PRIMARY KEY (sku, warehouse_id, order_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS reserves;

DROP TABLE IF EXISTS stocks;

DROP TABLE IF EXISTS order_items;

DROP TABLE IF EXISTS orders;
-- +goose StatementEnd

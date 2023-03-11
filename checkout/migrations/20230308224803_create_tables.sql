-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS carts (
    user_id bigint,
    sku integer,
    count smallint NOT NULL,
    PRIMARY KEY (user_id, sku)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS carts;
-- +goose StatementEnd

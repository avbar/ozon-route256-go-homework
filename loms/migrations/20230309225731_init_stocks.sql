-- +goose Up
-- +goose StatementBegin
INSERT INTO stocks (sku, warehouse_id, count) VALUES (6245113, 1, 10);
INSERT INTO stocks (sku, warehouse_id, count) VALUES (6245113, 2, 5);
INSERT INTO stocks (sku, warehouse_id, count) VALUES (6245113, 3, 7);
INSERT INTO stocks (sku, warehouse_id, count) VALUES (6966051, 2, 1);
INSERT INTO stocks (sku, warehouse_id, count) VALUES (6966051, 3, 3);
INSERT INTO stocks (sku, warehouse_id, count) VALUES (6967749, 1, 20);
INSERT INTO stocks (sku, warehouse_id, count) VALUES (6967749, 2, 9);
INSERT INTO stocks (sku, warehouse_id, count) VALUES (6967749, 3, 17);
INSERT INTO stocks (sku, warehouse_id, count) VALUES (7277168, 1, 8);
INSERT INTO stocks (sku, warehouse_id, count) VALUES (7895903, 2, 15);
INSERT INTO stocks (sku, warehouse_id, count) VALUES (8748527, 2, 99);
INSERT INTO stocks (sku, warehouse_id, count) VALUES (18247421, 3, 2);
INSERT INTO stocks (sku, warehouse_id, count) VALUES (19045918, 3, 6);
INSERT INTO stocks (sku, warehouse_id, count) VALUES (19065113, 1, 11);
INSERT INTO stocks (sku, warehouse_id, count) VALUES (19366373, 1, 30);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM stocks;
-- +goose StatementEnd

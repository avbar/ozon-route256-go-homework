package schema

import "time"

type Order struct {
	OrderID   int64     `db:"order_id"`
	Status    string    `db:"status"`
	User      int64     `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type OrderItem struct {
	OrderID int64  `db:"order_id"`
	SKU     uint32 `db:"sku"`
	Count   uint16 `db:"count"`
}

type Stock struct {
	SKU         uint32 `db:"sku"`
	WarehouseID int64  `db:"warehouse_id"`
	Count       uint64 `db:"count"`
}

type Reserve struct {
	SKU         uint32 `db:"sku"`
	WarehouseID int64  `db:"warehouse_id"`
	OrderID     int64  `db:"order_id"`
	Count       uint64 `db:"count"`
}

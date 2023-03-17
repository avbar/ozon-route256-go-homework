package schema

type CartItem struct {
	User  int64  `db:"user_id"`
	SKU   uint32 `db:"sku"`
	Count uint16 `db:"count"`
}

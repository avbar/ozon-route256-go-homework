package domain

type OrderID int64

type OrderItem struct {
	SKU   uint32
	Count uint16
}

type Order struct {
	Status string
	User   int64
	Items  []OrderItem
}

type Stock struct {
	WarehouseID int64
	Count       uint64
}

type Model struct {
}

func New() *Model {
	return &Model{}
}

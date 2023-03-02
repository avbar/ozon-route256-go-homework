package domain

import "context"

type CartItem struct {
	SKU   uint32
	Count uint16
	Name  string
	Price uint32
}

type Cart struct {
	Items      []CartItem
	TotalPrice uint32
}

type Stock struct {
	WarehouseID int64
	Count       uint64
}

type OrderItem struct {
	SKU   uint32
	Count uint16
}

type OrderID int64

type Product struct {
	Name  string
	Price uint32
}

type LOMSClient interface {
	Stocks(ctx context.Context, sku uint32) ([]Stock, error)
	CreateOrder(ctx context.Context, user int64, items []OrderItem) (OrderID, error)
}

type ProductClient interface {
	GetProduct(ctx context.Context, token string, sku uint32) (Product, error)
}

type Model struct {
	lomsClient    LOMSClient
	productClient ProductClient
}

func New(lomsClient LOMSClient, productClient ProductClient) *Model {
	return &Model{
		lomsClient:    lomsClient,
		productClient: productClient,
	}
}

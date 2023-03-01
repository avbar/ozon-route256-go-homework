package domain

import "context"

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

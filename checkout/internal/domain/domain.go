package domain

//go:generate sh -c "mkdir -p mocks && rm -rf mocks/loms_client_minimock.go mocks/product_client_minimock.go mocks/checkout_repository_minimock.go"
//go:generate minimock -i LOMSClient -o ./mocks/ -s "_minimock.go"
//go:generate minimock -i ProductClient -o ./mocks/ -s "_minimock.go"
//go:generate minimock -i CheckoutRepository -o ./mocks/ -s "_minimock.go"

import "context"

type CartItem struct {
	SKU   uint32
	Count uint16
}

type CartItemDetail struct {
	SKU   uint32
	Count uint16
	Name  string
	Price uint32
}

type Cart struct {
	Items      []CartItemDetail
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
	GetProducts(ctx context.Context, token string, skus []uint32) (map[uint32]Product, error)
}

type Model struct {
	lomsClient         LOMSClient
	productClient      ProductClient
	checkoutRepository CheckoutRepository
	transactionManager TransactionManager
}

func New(lomsClient LOMSClient, productClient ProductClient, checkoutRepository CheckoutRepository, transactionManager TransactionManager) *Model {
	return &Model{
		lomsClient:         lomsClient,
		productClient:      productClient,
		checkoutRepository: checkoutRepository,
		transactionManager: transactionManager,
	}
}

type CheckoutRepository interface {
	AddToCart(ctx context.Context, user int64, sku uint32, count uint16) error
	DeleteFromCart(ctx context.Context, user int64, sku uint32, count uint16) error
	ListCart(ctx context.Context, user int64) ([]CartItem, error)
}

type TransactionManager interface {
	RunRepeatableRead(ctx context.Context, f func(ctxTX context.Context) error) error
}

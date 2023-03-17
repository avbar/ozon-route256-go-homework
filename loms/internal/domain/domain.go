package domain

import (
	"context"
	"time"
)

type OrderID int64

const (
	OrderStatusNew             = "new"
	OrderStatusAwaitingPayment = "awaiting_payment"
	OrderStatusFailed          = "failed"
	OrderStatusPayed           = "payed"
	OrderStatusCancelled       = "cancelled"
)

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
	lomsRepository     LOMSRepository
	transactionManager TransactionManager
}

func New(lomsRepository LOMSRepository, transactionManager TransactionManager) *Model {
	m := &Model{
		lomsRepository:     lomsRepository,
		transactionManager: transactionManager,
	}

	go m.runOrderTimer(context.Background())

	return m
}

type LOMSRepository interface {
	CreateOrder(ctx context.Context, user int64, items []OrderItem) (OrderID, error)
	ListOrder(ctx context.Context, orderID OrderID) (Order, error)
	Stocks(ctx context.Context, sku uint32) ([]Stock, error)
	GetStatus(ctx context.Context, orderID OrderID) (string, error)
	ChangeStatus(ctx context.Context, orderID OrderID, status string) error
	CreateReserve(ctx context.Context, orderID OrderID, items []OrderItem) error
	CancelReserve(ctx context.Context, orderID OrderID) error
	DeleteReserve(ctx context.Context, orderID OrderID) error
	GetOldOrders(ctx context.Context, createdBefore time.Time) ([]OrderID, error)
}

type TransactionManager interface {
	RunRepeatableRead(ctx context.Context, f func(ctxTX context.Context) error) error
}

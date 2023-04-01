package domain

//go:generate sh -c "mkdir -p mocks && rm -rf mocks/loms_repository_minimock.go mocks/order_sender_minimock.go"
//go:generate minimock -i LOMSRepository -o ./mocks/ -s "_minimock.go"
//go:generate minimock -i OrderSender -o ./mocks/ -s "_minimock.go"

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

type OrdersOutbox struct {
	OrderID OrderID
	Status  string
}

type Model struct {
	lomsRepository     LOMSRepository
	transactionManager TransactionManager
	orderSender        OrderSender
}

func New(lomsRepository LOMSRepository, transactionManager TransactionManager, orderSender OrderSender) *Model {
	m := &Model{
		lomsRepository:     lomsRepository,
		transactionManager: transactionManager,
		orderSender:        orderSender,
	}

	ctx := context.Background()

	// Запускаем обработчик, отправляющий статусы заказов из Outbox в Kafka
	// В случае успешной отправки удаляем заказ из Outbox
	orderSender.AddSuccessHandler(ctx, func(orderID int64, status string) {
		m.lomsRepository.DeleteOrderFromOutbox(ctx, OrderID(orderID), status)
	})
	go m.runOutboxProcessor(ctx)

	// Обработка отмены заказов по тайм-ауту
	go m.runOrderTimer(ctx)

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

	SaveOrderToOutbox(ctx context.Context, orderID OrderID, status string) error
	DeleteOrderFromOutbox(ctx context.Context, orderID OrderID, status string) error
	GetOrdersFromOutbox(ctx context.Context) ([]OrdersOutbox, error)
}

type TransactionManager interface {
	RunRepeatableRead(ctx context.Context, f func(ctxTX context.Context) error) error
}

type Handler func(orderID int64, status string)

type OrderSender interface {
	SendOrderStatus(ctx context.Context, orderID int64, status string)
	AddSuccessHandler(ctx context.Context, onSuccess func(orderID int64, status string))
	AddErrorHandler(ctx context.Context, onError func(orderID int64, status string))
}

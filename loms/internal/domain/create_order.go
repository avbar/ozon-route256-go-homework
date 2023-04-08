package domain

import (
	"context"
)

func (m *Model) CreateOrder(ctx context.Context, user int64, items []OrderItem) (OrderID, error) {
	orderID, err := m.lomsRepository.CreateOrder(ctx, user, items)
	if err != nil {
		return 0, err
	}
	err = m.lomsRepository.SaveOrderToOutbox(ctx, orderID, OrderStatusNew)
	if err != nil {
		return 0, err
	}

	err = m.transactionManager.RunRepeatableRead(ctx, func(ctxTX context.Context) error {
		return m.lomsRepository.CreateReserve(ctx, orderID, items)
	})

	if err != nil {
		err = m.ChangeStatus(ctx, orderID, OrderStatusFailed)
	} else {
		err = m.ChangeStatus(ctx, orderID, OrderStatusAwaitingPayment)
	}

	return orderID, err
}

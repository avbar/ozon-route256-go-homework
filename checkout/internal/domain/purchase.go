package domain

import (
	"context"

	"github.com/pkg/errors"
)

func (m *Model) Purchase(ctx context.Context, user int64) (OrderID, error) {
	CartItems, err := m.checkoutRepository.ListCart(ctx, user)
	if err != nil {
		return 0, err
	}

	orderItems := make([]OrderItem, 0, len(CartItems))
	for _, cartItem := range CartItems {
		orderItems = append(orderItems, OrderItem{
			SKU:   cartItem.SKU,
			Count: cartItem.Count,
		})
	}

	var orderID OrderID
	err = m.transactionManager.RunRepeatableRead(ctx, func(ctxTX context.Context) error {
		orderID, err = m.lomsClient.CreateOrder(ctx, user, orderItems)
		if err != nil {
			return errors.WithMessage(err, "creating order")
		}

		for _, cartItem := range CartItems {
			err = m.checkoutRepository.DeleteFromCart(ctx, user, cartItem.SKU, cartItem.Count)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return 0, err
	}

	return orderID, nil
}

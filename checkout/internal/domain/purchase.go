package domain

import (
	"context"

	"github.com/pkg/errors"
)

func (m *Model) Purchase(ctx context.Context, user int64) (OrderID, error) {
	var orderItems = []OrderItem{
		{
			SKU:   777,
			Count: 7,
		},
		{
			SKU:   888,
			Count: 8,
		},
	}

	orderID, err := m.lomsClient.CreateOrder(ctx, user, orderItems)
	if err != nil {
		return orderID, errors.WithMessage(err, "checking stocks")
	}

	return orderID, nil
}

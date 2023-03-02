package domain

import (
	"context"
)

func (m *Model) ListOrder(ctx context.Context, orderID OrderID) (Order, error) {
	return Order{
		Status: "new",
		User:   123,
		Items: []OrderItem{
			{
				SKU:   456,
				Count: 2,
			},
		},
	}, nil
}

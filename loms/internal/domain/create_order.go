package domain

import (
	"context"
)

func (m *Model) CreateOrder(ctx context.Context, user int64, items []OrderItem) (OrderID, error) {
	return 12345, nil
}

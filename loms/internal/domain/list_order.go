package domain

import (
	"context"
)

func (m *Model) ListOrder(ctx context.Context, orderID OrderID) (Order, error) {
	return m.lomsRepository.ListOrder(ctx, orderID)
}

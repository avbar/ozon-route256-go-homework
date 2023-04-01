package domain

import (
	"context"
)

func (m *Model) OrderPayed(ctx context.Context, orderID OrderID) error {
	err := m.lomsRepository.DeleteReserve(ctx, orderID)
	if err != nil {
		return err
	}

	return m.ChangeStatus(ctx, orderID, OrderStatusPayed)
}

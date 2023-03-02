package domain

import (
	"context"
)

func (m *Model) Stocks(ctx context.Context, sku uint32) ([]Stock, error) {
	return []Stock{
		{
			WarehouseID: 123,
			Count:       5,
		},
		{
			WarehouseID: 456,
			Count:       3,
		},
	}, nil
}

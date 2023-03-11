package domain

import (
	"context"
)

func (m *Model) Stocks(ctx context.Context, sku uint32) ([]Stock, error) {
	return m.lomsRepository.Stocks(ctx, sku)
}

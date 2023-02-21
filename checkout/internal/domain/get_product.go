package domain

import (
	"context"

	"github.com/pkg/errors"
)

type Product struct {
	Name  string
	Price uint32
}

func (m *Model) GetProduct(ctx context.Context, token string, sku uint32) (Product, error) {
	product, err := m.productClient.GetProduct(ctx, token, sku)
	if err != nil {
		return product, errors.WithMessage(err, "getting product")
	}

	return product, nil
}

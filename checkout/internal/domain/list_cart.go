package domain

import (
	"context"
	"route256/checkout/internal/config"

	"github.com/pkg/errors"
)

func (m *Model) ListCart(ctx context.Context, user int64) (Cart, error) {
	token := config.ConfigData.Token
	var skus = []uint32{
		6245113,
		6966051,
		6967749,
	}

	var cart Cart

	for _, sku := range skus {
		product, err := m.productClient.GetProduct(ctx, token, sku)
		if err != nil {
			return cart, errors.WithMessage(err, "listing cart")
		}

		cart.Items = append(cart.Items,
			CartItem{
				SKU:   sku,
				Count: 1,
				Name:  product.Name,
				Price: product.Price,
			},
		)
		cart.TotalPrice += product.Price
	}

	return cart, nil
}

package domain

import (
	"context"
	"route256/checkout/internal/config"

	"github.com/pkg/errors"
)

func (m *Model) ListCart(ctx context.Context, user int64) (Cart, error) {
	token := config.ConfigData.Token
	var cart Cart

	CartItems, err := m.checkoutRepository.ListCart(ctx, user)
	if err != nil {
		return cart, nil
	}

	for _, cartItem := range CartItems {
		product, err := m.productClient.GetProduct(ctx, token, cartItem.SKU)
		if err != nil {
			return cart, errors.WithMessage(err, "listing cart")
		}

		cart.Items = append(cart.Items,
			CartItemDetail{
				SKU:   cartItem.SKU,
				Count: cartItem.Count,
				Name:  product.Name,
				Price: product.Price,
			},
		)
		cart.TotalPrice += product.Price
	}

	return cart, nil
}

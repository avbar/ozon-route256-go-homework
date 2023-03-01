package domain

import (
	"context"
	"route256/checkout/internal/config"

	"github.com/pkg/errors"
)

type CartItem struct {
	SKU   uint32
	Count uint16
	Name  string
	Price uint32
}

type Cart struct {
	Items      []CartItem
	TotalPrice uint32
}

func (m *Model) ListCart(ctx context.Context, user int64) (Cart, error) {
	token := config.ConfigData.Token
	var skus = []uint32{
		6245113,
		6966051,
		6967749,
	}

	var cartItems []CartItem
	var cart Cart

	for _, sku := range skus {
		product, err := m.productClient.GetProduct(ctx, token, sku)
		if err != nil {
			return cart, errors.WithMessage(err, "listing cart")
		}

		cartItems = append(cartItems,
			CartItem{
				SKU:   sku,
				Count: 1,
				Name:  product.Name,
				Price: product.Price,
			},
		)
	}

	for _, cartItem := range cartItems {
		cart.Items = append(cart.Items, cartItem)
		cart.TotalPrice += cartItem.Price
	}

	return cart, nil
}

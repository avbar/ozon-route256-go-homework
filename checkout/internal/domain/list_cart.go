package domain

import (
	"context"
	"route256/checkout/internal/config"

	"github.com/pkg/errors"
)

func (m *Model) ListCart(ctx context.Context, user int64) (Cart, error) {
	token := config.ConfigData.Token
	var cart Cart

	// Читаем из БД список товаров корзины
	CartItems, err := m.checkoutRepository.ListCart(ctx, user)
	if err != nil {
		return cart, nil
	}

	// Собираем все sku
	skus := make([]uint32, 0, len(CartItems))
	for _, cartItem := range CartItems {
		skus = append(skus, cartItem.SKU)
	}

	// Получаем описания товаров из ProductService по списку sku
	products, err := m.productClient.GetProducts(ctx, token, skus)
	if err != nil {
		return cart, errors.WithMessage(err, "listing cart")
	}

	// Формируем список товаров корзины с описаниями
	for _, cartItem := range CartItems {
		cart.Items = append(cart.Items,
			CartItemDetail{
				SKU:   cartItem.SKU,
				Count: cartItem.Count,
				Name:  products[cartItem.SKU].Name,
				Price: products[cartItem.SKU].Price,
			},
		)
		cart.TotalPrice += products[cartItem.SKU].Price
	}

	return cart, nil
}

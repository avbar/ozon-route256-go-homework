package checkout

import (
	"context"
	desc "route256/checkout/pkg/checkout_v1"
	"route256/libs/logger"

	"go.uber.org/zap"
)

func (i *Implementation) ListCart(ctx context.Context, req *desc.ListCartRequest) (*desc.ListCartResponse, error) {
	logger.Info("listCart", zap.Any("request", req))

	cart, err := i.businessLogic.ListCart(ctx, req.GetUser())
	if err != nil {
		return nil, err
	}

	items := make([]*desc.CartItem, 0, len(cart.Items))

	for _, item := range cart.Items {
		items = append(items, &desc.CartItem{
			Sku:   item.SKU,
			Count: uint32(item.Count),
			Name:  item.Name,
			Price: item.Price,
		})
	}

	return &desc.ListCartResponse{
		Items:      items,
		TotalPrice: cart.TotalPrice,
	}, nil
}

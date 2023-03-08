package checkout

import (
	"context"
	"log"
	desc "route256/checkout/pkg/checkout_v1"
)

func (i *Implementation) ListCart(ctx context.Context, req *desc.ListCartRequest) (*desc.ListCartResponse, error) {
	log.Printf("listCart: %+v", req)

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

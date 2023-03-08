package loms

import (
	"context"
	"log"
	"route256/loms/internal/domain"
	desc "route256/loms/pkg/loms_v1"
)

func (i *Implementation) ListOrder(ctx context.Context, req *desc.ListOrderRequest) (*desc.ListOrderResponse, error) {
	log.Printf("listOrder: %+v", req)

	order, err := i.businessLogic.ListOrder(ctx, domain.OrderID(req.GetOrderId()))
	if err != nil {
		return nil, err
	}

	items := make([]*desc.OrderItem, 0, len(order.Items))
	for _, item := range order.Items {
		items = append(items, &desc.OrderItem{
			Sku:   item.SKU,
			Count: uint32(item.Count),
		})
	}

	return &desc.ListOrderResponse{
		Status: desc.OrderStatus(desc.OrderStatus_value[order.Status]),
		User:   order.User,
		Items:  items,
	}, nil
}

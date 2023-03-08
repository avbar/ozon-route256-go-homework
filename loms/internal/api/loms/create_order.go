package loms

import (
	"context"
	"log"
	"route256/loms/internal/domain"
	desc "route256/loms/pkg/loms_v1"
)

func (i *Implementation) CreateOrder(ctx context.Context, req *desc.CreateOrderRequest) (*desc.CreateOrderResponse, error) {
	log.Printf("createOrder: %+v", req)

	reqItems := req.GetItems()
	items := make([]domain.OrderItem, 0, len(reqItems))
	for _, item := range reqItems {
		items = append(items, domain.OrderItem{
			SKU:   item.GetSku(),
			Count: uint16(item.GetCount()),
		})
	}

	orderID, err := i.businessLogic.CreateOrder(ctx, req.GetUser(), items)
	if err != nil {
		return nil, err
	}

	return &desc.CreateOrderResponse{
		OrderId: int64(orderID),
	}, nil
}

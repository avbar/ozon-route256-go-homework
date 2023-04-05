package loms

import (
	"context"
	"route256/libs/logger"
	"route256/loms/internal/domain"
	desc "route256/loms/pkg/loms_v1"

	"go.uber.org/zap"
)

func (i *Implementation) ListOrder(ctx context.Context, req *desc.ListOrderRequest) (*desc.ListOrderResponse, error) {
	logger.Info("listOrder", zap.Any("request", req))

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

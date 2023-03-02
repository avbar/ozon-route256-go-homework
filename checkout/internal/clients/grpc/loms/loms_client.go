package loms

import (
	"context"
	"route256/checkout/internal/domain"
	lomsAPI "route256/loms/pkg/loms_v1"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

var _ Client = (*client)(nil)

type Client interface {
	Stocks(ctx context.Context, sku uint32) ([]domain.Stock, error)
	CreateOrder(ctx context.Context, user int64, items []domain.OrderItem) (domain.OrderID, error)
}

type client struct {
	lomsClient lomsAPI.LOMSClient
}

func NewClient(cc *grpc.ClientConn) *client {
	return &client{
		lomsClient: lomsAPI.NewLOMSClient(cc),
	}
}

func (c *client) Stocks(ctx context.Context, sku uint32) ([]domain.Stock, error) {
	res, err := c.lomsClient.Stocks(ctx, &lomsAPI.StocksRequest{Sku: sku})
	if err != nil {
		return nil, errors.Wrap(err, "calling Stocks")
	}

	stocks := make([]domain.Stock, 0, len(res.Stocks))
	for _, stock := range res.Stocks {
		stocks = append(stocks, domain.Stock{
			WarehouseID: stock.GetWarehouseId(),
			Count:       stock.GetCount(),
		})
	}

	return stocks, nil
}

func (c *client) CreateOrder(ctx context.Context, user int64, items []domain.OrderItem) (domain.OrderID, error) {
	orderItems := make([]*lomsAPI.OrderItem, 0, len(items))
	for _, item := range items {
		orderItems = append(orderItems, &lomsAPI.OrderItem{
			Sku:   item.SKU,
			Count: uint32(item.Count),
		})
	}

	res, err := c.lomsClient.CreateOrder(ctx, &lomsAPI.CreateOrderRequest{
		User:  user,
		Items: orderItems,
	})
	if err != nil {
		return 0, errors.Wrap(err, "calling CreateOrder")
	}

	return domain.OrderID(res.GetOrderId()), nil
}

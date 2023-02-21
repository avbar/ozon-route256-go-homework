package loms

import (
	"context"
	"route256/checkout/internal/domain"
	"route256/libs/clntwrapper"

	"github.com/pkg/errors"
)

type Client struct {
	url string

	urlStocks      string
	urlCreateOrder string
}

func New(url string) *Client {
	return &Client{
		url: url,

		urlStocks:      url + "/stocks",
		urlCreateOrder: url + "/createOrder",
	}
}

type StocksRequest struct {
	SKU uint32 `json:"sku"`
}

type StocksItem struct {
	WarehouseID int64  `json:"warehouseID"`
	Count       uint64 `json:"count"`
}

type StocksResponse struct {
	Stocks []StocksItem `json:"stocks"`
}

func (c *Client) Stocks(ctx context.Context, sku uint32) ([]domain.Stock, error) {
	request := StocksRequest{SKU: sku}

	wrapper := clntwrapper.New[StocksRequest, StocksResponse](c.urlStocks)
	response, err := wrapper.Handle(ctx, request)
	if err != nil {
		return nil, errors.Wrap(err, "calling stocks")
	}

	stocks := make([]domain.Stock, 0, len(response.Stocks))
	for _, stock := range response.Stocks {
		stocks = append(stocks, domain.Stock{
			WarehouseID: stock.WarehouseID,
			Count:       stock.Count,
		})
	}

	return stocks, nil
}

type OrderItem struct {
	SKU   uint32 `json:"sku"`
	Count uint16 `json:"count"`
}

type CreateOrderRequest struct {
	User  int64       `json:"user"`
	Items []OrderItem `json:"items"`
}

type CreateOrderResponse struct {
	OrderID int64 `json:"orderID"`
}

func (c *Client) CreateOrder(ctx context.Context, user int64, items []domain.OrderItem) (domain.OrderID, error) {
	var orderItems []OrderItem
	for _, item := range items {
		orderItems = append(orderItems, OrderItem{SKU: item.SKU, Count: item.Count})
	}
	request := CreateOrderRequest{User: user, Items: orderItems}

	wrapper := clntwrapper.New[CreateOrderRequest, CreateOrderResponse](c.urlCreateOrder)
	response, err := wrapper.Handle(ctx, request)
	if err != nil {
		return 0, errors.Wrap(err, "calling createOrder")
	}

	orderID := domain.OrderID(response.OrderID)

	return orderID, nil
}

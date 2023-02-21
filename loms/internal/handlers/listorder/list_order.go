package listorder

import (
	"context"
	"errors"
	"log"
)

type Request struct {
	OrderID int64 `json:"orderID"`
}

type OrderItem struct {
	SKU   uint32 `json:"sku"`
	Count uint16 `json:"count"`
}

type Response struct {
	Status string      `json:"status"`
	User   int64       `json:"user"`
	Items  []OrderItem `json:"items"`
}

var (
	ErrEmptyOrderID = errors.New("empty order ID")
)

func (r Request) Validate() error {
	if r.OrderID == 0 {
		return ErrEmptyOrderID
	}
	return nil
}

type Handler struct{}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) Handle(ctx context.Context, request Request) (Response, error) {
	log.Printf("listOrder: %+v", request)
	return Response{
		Status: "new",
		User:   123,
		Items: []OrderItem{
			{
				SKU:   456,
				Count: 2,
			},
		},
	}, nil
}

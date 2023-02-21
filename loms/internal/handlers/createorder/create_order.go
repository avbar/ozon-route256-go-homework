package createorder

import (
	"context"
	"errors"
	"log"
)

type OrderItem struct {
	SKU   uint32 `json:"sku"`
	Count uint16 `json:"count"`
}

type Request struct {
	User  int64       `json:"user"`
	Items []OrderItem `json:"items"`
}

type Response struct {
	OrderID int64 `json:"orderID"`
}

var (
	ErrEmptyUser = errors.New("empty user")
)

func (r Request) Validate() error {
	if r.User == 0 {
		return ErrEmptyUser
	}
	return nil
}

type Handler struct{}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) Handle(ctx context.Context, request Request) (Response, error) {
	log.Printf("createOrder: %+v", request)
	return Response{
		OrderID: 12345,
	}, nil
}

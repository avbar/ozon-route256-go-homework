package orderpayed

import (
	"context"
	"errors"
	"log"
)

type Request struct {
	OrderID int64 `json:"orderID"`
}

type Response struct{}

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
	log.Printf("orderPayed: %+v", request)
	return Response{}, nil
}

package listcart

import (
	"context"
	"errors"
	"log"
	"route256/checkout/internal/domain"
)

type Handler struct {
	businessLogic *domain.Model
}

func New(businessLogic *domain.Model) *Handler {
	return &Handler{
		businessLogic: businessLogic,
	}
}

type Request struct {
	User int64 `json:"user"`
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

type CartItem struct {
	SKU   uint32 `json:"sku"`
	Count uint16 `json:"count"`
	Name  string `json:"name"`
	Price uint32 `json:"price"`
}

type Response struct {
	Items      []CartItem `json:"items"`
	TotalPrice uint32     `json:"totalPrice"`
}

func (h *Handler) Handle(ctx context.Context, req Request) (Response, error) {
	log.Printf("listCart: %+v", req)

	var response Response

	cart, err := h.businessLogic.ListCart(ctx, req.User)
	if err != nil {
		return response, err
	}

	for _, item := range cart.Items {
		response.Items = append(response.Items, CartItem{
			SKU:   item.SKU,
			Count: item.Count,
			Name:  item.Name,
			Price: item.Price,
		})
	}
	response.TotalPrice = cart.TotalPrice

	return response, nil
}

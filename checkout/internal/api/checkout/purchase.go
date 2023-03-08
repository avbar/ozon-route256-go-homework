package checkout

import (
	"context"
	"log"
	desc "route256/checkout/pkg/checkout_v1"
)

func (i *Implementation) Purchase(ctx context.Context, req *desc.PurchaseRequest) (*desc.PurchaseResponse, error) {
	log.Printf("purchase: %+v", req)

	orderID, err := i.businessLogic.Purchase(ctx, req.GetUser())
	if err != nil {
		return nil, err
	}

	return &desc.PurchaseResponse{
		OrderId: int64(orderID),
	}, nil
}

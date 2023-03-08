package loms

import (
	"context"
	"log"
	"route256/loms/internal/domain"
	desc "route256/loms/pkg/loms_v1"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (i *Implementation) CancelOrder(ctx context.Context, req *desc.CancelOrderRequest) (*emptypb.Empty, error) {
	log.Printf("cancelOrder: %+v", req)

	err := i.businessLogic.CancelOrder(ctx, domain.OrderID(req.GetOrderId()))
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

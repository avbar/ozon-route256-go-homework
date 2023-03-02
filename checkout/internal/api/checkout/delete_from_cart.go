package checkout

import (
	"context"
	"log"
	desc "route256/checkout/pkg/checkout_v1"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (i *Implementation) DeleteFromCart(ctx context.Context, req *desc.DeleteFromCartRequest) (*emptypb.Empty, error) {
	log.Printf("deleteFromCart: %+v", req)

	err := i.businessLogic.DeleteFromCart(ctx, req.GetUser(), req.GetSku(), uint16(req.GetCount()))
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

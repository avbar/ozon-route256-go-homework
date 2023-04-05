package checkout

import (
	"context"
	desc "route256/checkout/pkg/checkout_v1"
	"route256/libs/logger"

	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (i *Implementation) AddToCart(ctx context.Context, req *desc.AddToCartRequest) (*emptypb.Empty, error) {
	logger.Info("addToCart", zap.Any("request", req))

	err := i.businessLogic.AddToCart(ctx, req.GetUser(), req.GetSku(), uint16(req.GetCount()))
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

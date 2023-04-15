package loms

import (
	"context"
	"route256/libs/logger"
	"route256/loms/internal/domain"
	desc "route256/loms/pkg/loms_v1"

	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (i *Implementation) OrderPayed(ctx context.Context, req *desc.OrderPayedRequest) (*emptypb.Empty, error) {
	logger.Info("orderPayed", zap.Any("request", req))

	err := i.businessLogic.OrderPayed(ctx, domain.OrderID(req.GetOrderId()))
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

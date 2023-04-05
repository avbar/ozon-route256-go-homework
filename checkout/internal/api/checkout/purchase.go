package checkout

import (
	"context"
	desc "route256/checkout/pkg/checkout_v1"
	"route256/libs/logger"

	"go.uber.org/zap"
)

func (i *Implementation) Purchase(ctx context.Context, req *desc.PurchaseRequest) (*desc.PurchaseResponse, error) {
	logger.Info("purchase", zap.Any("request", req))

	orderID, err := i.businessLogic.Purchase(ctx, req.GetUser())
	if err != nil {
		return nil, err
	}

	return &desc.PurchaseResponse{
		OrderId: int64(orderID),
	}, nil
}

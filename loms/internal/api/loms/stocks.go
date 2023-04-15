package loms

import (
	"context"
	"route256/libs/logger"
	desc "route256/loms/pkg/loms_v1"

	"go.uber.org/zap"
)

func (i *Implementation) Stocks(ctx context.Context, req *desc.StocksRequest) (*desc.StocksResponse, error) {
	logger.Info("stocks", zap.Any("request", req))

	stocks, err := i.businessLogic.Stocks(ctx, req.GetSku())
	if err != nil {
		return nil, err
	}

	resStocks := make([]*desc.Stock, 0, len(stocks))
	for _, stock := range stocks {
		resStocks = append(resStocks, &desc.Stock{
			WarehouseId: stock.WarehouseID,
			Count:       stock.Count,
		})
	}

	return &desc.StocksResponse{
		Stocks: resStocks,
	}, nil
}

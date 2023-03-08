package loms

import (
	"context"
	"log"
	desc "route256/loms/pkg/loms_v1"
)

func (i *Implementation) Stocks(ctx context.Context, req *desc.StocksRequest) (*desc.StocksResponse, error) {
	log.Printf("stocks: %+v", req)

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

package repository

import (
	"context"
	"route256/loms/internal/domain"
	"route256/loms/internal/repository/schema"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
)

func (r *LOMSRepo) Stocks(ctx context.Context, sku uint32) ([]domain.Stock, error) {
	db := r.QueryEngineProvider.GetQueryEngine(ctx)

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query := builder.Select(stocksColumns...).
		From(stocksTable).
		Where(sq.Eq{"sku": sku})

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var stocks []schema.Stock
	if err := pgxscan.Select(ctx, db, &stocks, rawQuery, args...); err != nil {
		return nil, err
	}

	return bindSchemaStocksToDomainStocks(stocks), nil
}

func bindSchemaStocksToDomainStocks(stocks []schema.Stock) []domain.Stock {
	result := make([]domain.Stock, 0, len(stocks))
	for _, stock := range stocks {
		result = append(result, domain.Stock{
			WarehouseID: stock.WarehouseID,
			Count:       stock.Count,
		})
	}
	return result
}

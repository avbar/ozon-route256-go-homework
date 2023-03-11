package repository

import (
	"context"
	"route256/loms/internal/domain"
	"route256/loms/internal/repository/schema"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
)

func (r *LOMSRepo) CancelReserve(ctx context.Context, orderID domain.OrderID) error {
	db := r.QueryEngineProvider.GetQueryEngine(ctx)

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	// get reserves
	query := builder.Select(reservesColumns...).
		From(reservesTable).
		Where(sq.Eq{"order_id": orderID})

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return err
	}

	var reserves []schema.Reserve
	if err := pgxscan.Select(ctx, db, &reserves, rawQuery, args...); err != nil {
		return err
	}

	for _, reserve := range reserves {
		// insert into stocks
		insertQuery := builder.Insert(stocksTable).
			Columns(stocksColumns...).
			Values(reserve.SKU, reserve.WarehouseID, reserve.Count).
			Suffix("ON CONFLICT ON CONSTRAINT stocks_pkey DO UPDATE SET count = stocks.count + excluded.count")

		rawQuery, args, err := insertQuery.ToSql()
		if err != nil {
			return err
		}

		_, err = db.Exec(ctx, rawQuery, args...)
		if err != nil {
			return err
		}
	}

	// delete from reserves
	err = r.DeleteReserve(ctx, orderID)
	if err != nil {
		return nil
	}

	return nil
}

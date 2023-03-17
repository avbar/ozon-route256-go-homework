package repository

import (
	"context"
	"route256/loms/internal/domain"

	sq "github.com/Masterminds/squirrel"
)

func (r *LOMSRepo) CreateReserve(ctx context.Context, orderID domain.OrderID, items []domain.OrderItem) error {
	db := r.QueryEngineProvider.GetQueryEngine(ctx)

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	for _, item := range items {
		// insert items into order
		insertQuery := builder.Insert(orderItemsTable).
			Columns(orderItemsColumns...).
			Values(orderID, item.SKU, item.Count)

		rawQuery, args, err := insertQuery.ToSql()
		if err != nil {
			return err
		}

		_, err = db.Exec(ctx, rawQuery, args...)
		if err != nil {
			return err
		}

		// check and update stocks
		restCount := item.Count

		stocks, err := r.Stocks(ctx, item.SKU)
		if err != nil {
			return err
		}

		for _, stock := range stocks {
			var count uint16
			if uint64(restCount) < stock.Count {
				count = restCount
			} else {
				count = uint16(stock.Count)
			}

			// insert into reserves
			insertQuery := builder.Insert(reservesTable).
				Columns(reservesColumns...).
				Values(item.SKU, stock.WarehouseID, orderID, count)

			rawQuery, args, err := insertQuery.ToSql()
			if err != nil {
				return err
			}

			_, err = db.Exec(ctx, rawQuery, args...)
			if err != nil {
				return err
			}

			// delete from stocks
			if uint64(count) == stock.Count {
				deleteQuery := builder.Delete(stocksTable).
					Where(sq.Eq{"sku": item.SKU, "warehouse_id": stock.WarehouseID})

				rawQuery, args, err = deleteQuery.ToSql()
				if err != nil {
					return err
				}
			} else {
				updateQuery := builder.Update(stocksTable).
					Set("count", stock.Count-uint64(count)).
					Where(sq.Eq{"sku": item.SKU, "warehouse_id": stock.WarehouseID})

				rawQuery, args, err = updateQuery.ToSql()
				if err != nil {
					return err
				}
			}

			_, err = db.Exec(ctx, rawQuery, args...)
			if err != nil {
				return err
			}

			restCount -= count
			if restCount <= 0 {
				break
			}
		}
	}

	return nil
}

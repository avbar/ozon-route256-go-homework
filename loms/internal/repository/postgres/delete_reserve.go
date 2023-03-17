package repository

import (
	"context"
	"route256/loms/internal/domain"

	sq "github.com/Masterminds/squirrel"
)

func (r *LOMSRepo) DeleteReserve(ctx context.Context, orderID domain.OrderID) error {
	db := r.QueryEngineProvider.GetQueryEngine(ctx)

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	deleteQuery := builder.Delete(reservesTable).
		Where(sq.Eq{"order_id": orderID})

	rawQuery, args, err := deleteQuery.ToSql()
	if err != nil {
		return err
	}

	_, err = db.Exec(ctx, rawQuery, args...)
	if err != nil {
		return err
	}

	return nil
}

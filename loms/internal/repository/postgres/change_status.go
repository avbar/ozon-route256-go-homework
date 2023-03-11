package repository

import (
	"context"
	"route256/loms/internal/domain"
	"time"

	sq "github.com/Masterminds/squirrel"
)

func (r *LOMSRepo) ChangeStatus(ctx context.Context, orderID domain.OrderID, status string) error {
	db := r.QueryEngineProvider.GetQueryEngine(ctx)

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query := builder.Update(ordersTable).
		Set("status", status).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"order_id": orderID})

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = db.Exec(ctx, rawQuery, args...)
	if err != nil {
		return err
	}

	return nil
}

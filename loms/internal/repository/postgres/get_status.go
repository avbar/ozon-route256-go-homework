package repository

import (
	"context"
	"route256/loms/internal/domain"

	sq "github.com/Masterminds/squirrel"
)

func (r *LOMSRepo) GetStatus(ctx context.Context, orderID domain.OrderID) (string, error) {
	db := r.QueryEngineProvider.GetQueryEngine(ctx)

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query := builder.Select("status").
		From(ordersTable).
		Where(sq.Eq{"order_id": orderID})

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return "", err
	}

	var status string
	err = db.QueryRow(ctx, rawQuery, args...).Scan(&status)
	if err != nil {
		return "", err
	}

	return status, nil
}

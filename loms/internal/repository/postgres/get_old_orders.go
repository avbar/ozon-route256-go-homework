package repository

import (
	"context"
	"route256/loms/internal/domain"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
)

func (r *LOMSRepo) GetOldOrders(ctx context.Context, createdBefore time.Time) ([]domain.OrderID, error) {
	db := r.QueryEngineProvider.GetQueryEngine(ctx)

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query := builder.Select("order_id").
		From(ordersTable).
		Where(sq.Or{sq.Eq{"status": domain.OrderStatusNew}, sq.Eq{"status": domain.OrderStatusAwaitingPayment}}).
		Where(sq.LtOrEq{"created_at": createdBefore})

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var orders []domain.OrderID
	if err := pgxscan.Select(ctx, db, &orders, rawQuery, args...); err != nil {
		return nil, err
	}

	return orders, nil
}

package repository

import (
	"context"
	"route256/loms/internal/domain"
	"time"

	sq "github.com/Masterminds/squirrel"
)

func (r *LOMSRepo) CreateOrder(ctx context.Context, user int64, items []domain.OrderItem) (domain.OrderID, error) {
	db := r.QueryEngineProvider.GetQueryEngine(ctx)

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query := builder.Insert(ordersTable).
		Columns("status", "user_id", "created_at", "updated_at").
		Values(domain.OrderStatusNew, user, time.Now(), time.Now()).
		Suffix("RETURNING order_id")

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return 0, err
	}

	var orderID domain.OrderID
	err = db.QueryRow(ctx, rawQuery, args...).Scan(&orderID)
	if err != nil {
		return 0, err
	}

	return orderID, nil
}

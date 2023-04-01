package repository

import (
	"context"
	"route256/loms/internal/domain"
	"route256/loms/internal/repository/schema"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
)

func (r *LOMSRepo) SaveOrderToOutbox(ctx context.Context, orderID domain.OrderID, status string) error {
	db := r.QueryEngineProvider.GetQueryEngine(ctx)

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query := builder.Insert(ordersOutboxTable).
		Columns("order_id", "status", "created_at").
		Values(orderID, status, time.Now())

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

func (r *LOMSRepo) DeleteOrderFromOutbox(ctx context.Context, orderID domain.OrderID, status string) error {
	db := r.QueryEngineProvider.GetQueryEngine(ctx)

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query := builder.Delete(ordersOutboxTable).
		Where(sq.Eq{"order_id": orderID, "status": status})

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

func (r *LOMSRepo) GetOrdersFromOutbox(ctx context.Context) ([]domain.OrdersOutbox, error) {
	db := r.QueryEngineProvider.GetQueryEngine(ctx)

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query := builder.Select(ordersOutboxColumns...).
		From(ordersOutboxTable).
		OrderBy("created_at")

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var outbox []schema.OrdersOutbox
	if err := pgxscan.Select(ctx, db, &outbox, rawQuery, args...); err != nil {
		return nil, err
	}

	var ordersOutbox []domain.OrdersOutbox
	for _, order := range outbox {
		ordersOutbox = append(ordersOutbox, domain.OrdersOutbox{
			OrderID: domain.OrderID(order.OrderID),
			Status:  order.Status,
		})
	}
	return ordersOutbox, nil
}

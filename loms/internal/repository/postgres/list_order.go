package repository

import (
	"context"
	"errors"
	"route256/loms/internal/domain"
	"route256/loms/internal/repository/schema"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
)

var (
	ErrWrongOrder = errors.New("wrong order")
)

func (r *LOMSRepo) ListOrder(ctx context.Context, orderID domain.OrderID) (domain.Order, error) {
	db := r.QueryEngineProvider.GetQueryEngine(ctx)

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query := builder.Select(ordersColumns...).
		From(ordersTable).
		Where(sq.Eq{"order_id": orderID})

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return domain.Order{}, err
	}

	var orders []schema.Order
	if err := pgxscan.Select(ctx, db, &orders, rawQuery, args...); err != nil {
		return domain.Order{}, err
	}
	if len(orders) == 0 {
		return domain.Order{}, ErrWrongOrder
	}
	order := orders[0]

	query = builder.Select(orderItemsColumns...).
		From(orderItemsTable).
		Where(sq.Eq{"order_id": orderID})

	rawQuery, args, err = query.ToSql()
	if err != nil {
		return domain.Order{}, err
	}

	var orderItems []schema.OrderItem
	if err := pgxscan.Select(ctx, db, &orderItems, rawQuery, args...); err != nil {
		return domain.Order{}, err
	}

	return bindSchemaOrderToDomainOrder(order, orderItems), nil
}

func bindSchemaOrderToDomainOrder(order schema.Order, orderItems []schema.OrderItem) domain.Order {
	var result domain.Order
	result.Status = order.Status
	result.User = order.User

	result.Items = make([]domain.OrderItem, 0, len(orderItems))
	for _, item := range orderItems {
		result.Items = append(result.Items, domain.OrderItem{
			SKU:   item.SKU,
			Count: item.Count,
		})
	}

	return result
}

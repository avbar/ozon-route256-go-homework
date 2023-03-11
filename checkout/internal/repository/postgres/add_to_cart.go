package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
)

func (r *CheckoutRepo) AddToCart(ctx context.Context, user int64, sku uint32, count uint16) error {
	db := r.QueryEngineProvider.GetQueryEngine(ctx)

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query := builder.Insert(cartsTable).
		Columns(cartsColumns...).
		Values(user, sku, count).
		Suffix("ON CONFLICT ON CONSTRAINT carts_pkey DO UPDATE SET count = carts.count + excluded.count")

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

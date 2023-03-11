package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
)

func (r *CheckoutRepo) DeleteFromCart(ctx context.Context, user int64, sku uint32, count uint16) error {
	db := r.QueryEngineProvider.GetQueryEngine(ctx)

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query := builder.Select("count").
		From(cartsTable).
		Where(sq.Eq{"user_id": user, "sku": sku})

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return err
	}

	var currCount uint16
	err = db.QueryRow(ctx, rawQuery, args...).Scan(&currCount)
	if err != nil {
		return err
	}

	if currCount <= count {
		// delete
		deleteQuery := builder.Delete(cartsTable).
			Where(sq.Eq{"user_id": user, "sku": sku})

		rawQuery, args, err = deleteQuery.ToSql()
		if err != nil {
			return err
		}
	} else {
		// update
		updateQuery := builder.Update(cartsTable).
			Set("count", currCount-count).
			Where(sq.Eq{"user_id": user, "sku": sku})

		rawQuery, args, err = updateQuery.ToSql()
		if err != nil {
			return err
		}
	}

	_, err = db.Exec(ctx, rawQuery, args...)
	if err != nil {
		return err
	}

	return nil
}

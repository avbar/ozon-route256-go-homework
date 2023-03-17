package repository

import (
	"context"
	"route256/checkout/internal/domain"
	"route256/checkout/internal/repository/schema"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
)

func (r *CheckoutRepo) ListCart(ctx context.Context, user int64) ([]domain.CartItem, error) {
	db := r.QueryEngineProvider.GetQueryEngine(ctx)

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query := builder.Select(cartsColumns...).
		From(cartsTable).
		Where(sq.Eq{"user_id": user})

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var cartItems []schema.CartItem
	if err := pgxscan.Select(ctx, db, &cartItems, rawQuery, args...); err != nil {
		return nil, err
	}

	return bindSchemaCartToDomainCart(cartItems), nil
}

func bindSchemaCartToDomainCart(cartItems []schema.CartItem) []domain.CartItem {
	result := make([]domain.CartItem, 0, len(cartItems))
	for _, cartItem := range cartItems {
		result = append(result, domain.CartItem{
			SKU:   cartItem.SKU,
			Count: cartItem.Count,
		})
	}
	return result
}

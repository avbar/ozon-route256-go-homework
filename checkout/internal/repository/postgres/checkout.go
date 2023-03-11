package repository

import (
	"route256/checkout/internal/domain"
	"route256/libs/postgres/transactor"
)

var _ domain.CheckoutRepository = (*CheckoutRepo)(nil)

type CheckoutRepo struct {
	transactor.QueryEngineProvider
}

func NewCheckoutRepo(provider transactor.QueryEngineProvider) *CheckoutRepo {
	return &CheckoutRepo{
		QueryEngineProvider: provider,
	}
}

var (
	cartsColumns = []string{"user_id", "sku", "count"}
)

const (
	cartsTable = "carts"
)

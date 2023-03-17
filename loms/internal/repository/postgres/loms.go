package repository

import (
	"route256/libs/postgres/transactor"
	"route256/loms/internal/domain"
)

var _ domain.LOMSRepository = (*LOMSRepo)(nil)

type LOMSRepo struct {
	transactor.QueryEngineProvider
}

func NewLOMSRepo(provider transactor.QueryEngineProvider) *LOMSRepo {
	return &LOMSRepo{
		QueryEngineProvider: provider,
	}
}

var (
	ordersColumns     = []string{"order_id", "status", "user_id", "created_at", "updated_at"}
	orderItemsColumns = []string{"order_id", "sku", "count"}
	stocksColumns     = []string{"sku", "warehouse_id", "count"}
	reservesColumns   = []string{"sku", "warehouse_id", "order_id", "count"}
)

const (
	ordersTable     = "orders"
	orderItemsTable = "order_items"
	stocksTable     = "stocks"
	reservesTable   = "reserves"
)

package domain_test

import (
	"route256/checkout/internal/domain"
	"route256/libs/postgres/transactor"

	"github.com/gojuno/minimock/v3"
)

type lomsClientMockFunc func(mc *minimock.Controller) domain.LOMSClient
type productClientMockFunc func(mc *minimock.Controller) domain.ProductClient
type checkoutRepositoryMockFunc func(mc *minimock.Controller) domain.CheckoutRepository
type dbMockFunc func(mc *minimock.Controller) transactor.DB

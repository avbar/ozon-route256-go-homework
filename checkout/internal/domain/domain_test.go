package domain_test

import (
	"route256/checkout/internal/domain"

	"github.com/gojuno/minimock/v3"
)

type lomsClientMockFunc func(mc *minimock.Controller) domain.LOMSClient
type productClientMockFunc func(mc *minimock.Controller) domain.ProductClient
type checkoutRepositoryMockFunc func(mc *minimock.Controller) domain.CheckoutRepository

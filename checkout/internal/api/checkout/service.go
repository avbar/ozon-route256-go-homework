package checkout

import (
	"route256/checkout/internal/domain"
	desc "route256/checkout/pkg/checkout_v1"
)

type Implementation struct {
	desc.UnimplementedCheckoutServer

	businessLogic *domain.Model
}

func NewCheckout(businessLogic *domain.Model) *Implementation {
	return &Implementation{
		businessLogic: businessLogic,
	}
}

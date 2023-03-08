package loms

import (
	"route256/loms/internal/domain"
	desc "route256/loms/pkg/loms_v1"
)

type Implementation struct {
	desc.UnimplementedLOMSServer

	businessLogic *domain.Model
}

func NewLOMS(businessLogic *domain.Model) *Implementation {
	return &Implementation{
		businessLogic: businessLogic,
	}
}

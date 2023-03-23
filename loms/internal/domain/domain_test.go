package domain_test

import (
	"route256/libs/postgres/transactor"
	"route256/loms/internal/domain"

	"github.com/gojuno/minimock/v3"
)

type lomsRepositoryMockFunc func(mc *minimock.Controller) domain.LOMSRepository
type dbMockFunc func(mc *minimock.Controller) transactor.DB

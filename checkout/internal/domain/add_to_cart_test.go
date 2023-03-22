package domain_test

import (
	"context"
	"route256/checkout/internal/domain"
	"route256/checkout/internal/domain/mocks"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestAddToCart(t *testing.T) {
	type args struct {
		ctx   context.Context
		user  int64
		sku   uint32
		count uint16
	}

	var (
		mc  = minimock.NewController(t)
		ctx = context.Background()

		user                     = gofakeit.Int64()
		sku                      = gofakeit.Uint32()
		countSufficient   uint16 = 10
		countInsufficient uint16 = 100

		repoErr = errors.New("repo add to cart error")
		lomsErr = errors.New("loms stocks error")

		lomsStocksRes []domain.Stock
	)

	for i := 1; i <= 5; i++ {
		lomsStocksRes = append(lomsStocksRes, domain.Stock{
			WarehouseID: gofakeit.Int64(),
			Count:       uint64(i),
		})
	}

	tests := []struct {
		name                   string
		args                   args
		err                    error
		checkoutRepositoryMock checkoutRepositoryMockFunc
		lomsClientMock         lomsClientMockFunc
	}{
		{
			name: "positive case",
			args: args{
				ctx:   ctx,
				user:  user,
				sku:   sku,
				count: countSufficient,
			},
			err: nil,
			checkoutRepositoryMock: func(mc *minimock.Controller) domain.CheckoutRepository {
				mock := mocks.NewCheckoutRepositoryMock(mc)
				mock.AddToCartMock.Expect(ctx, user, sku, countSufficient).Return(nil)
				return mock
			},
			lomsClientMock: func(mc *minimock.Controller) domain.LOMSClient {
				mock := mocks.NewLOMSClientMock(mc)
				mock.StocksMock.Expect(ctx, sku).Return(lomsStocksRes, nil)
				return mock
			},
		},

		{
			name: "insufficient stocks",
			args: args{
				ctx:   ctx,
				user:  user,
				sku:   sku,
				count: countInsufficient,
			},
			err: domain.ErrInsufficientStocks,
			checkoutRepositoryMock: func(mc *minimock.Controller) domain.CheckoutRepository {
				mock := mocks.NewCheckoutRepositoryMock(mc)
				mock.AddToCartMock.Expect(ctx, user, sku, countInsufficient).Return(nil)
				return mock
			},
			lomsClientMock: func(mc *minimock.Controller) domain.LOMSClient {
				mock := mocks.NewLOMSClientMock(mc)
				mock.StocksMock.Expect(ctx, sku).Return(lomsStocksRes, nil)
				return mock
			},
		},

		{
			name: "repo error",
			args: args{
				ctx:   ctx,
				user:  user,
				sku:   sku,
				count: countSufficient,
			},
			err: repoErr,
			checkoutRepositoryMock: func(mc *minimock.Controller) domain.CheckoutRepository {
				mock := mocks.NewCheckoutRepositoryMock(mc)
				mock.AddToCartMock.Expect(ctx, user, sku, countSufficient).Return(repoErr)
				return mock
			},
			lomsClientMock: func(mc *minimock.Controller) domain.LOMSClient {
				mock := mocks.NewLOMSClientMock(mc)
				mock.StocksMock.Expect(ctx, sku).Return(lomsStocksRes, nil)
				return mock
			},
		},

		{
			name: "loms service error",
			args: args{
				ctx:   ctx,
				user:  user,
				sku:   sku,
				count: countSufficient,
			},
			err: lomsErr,
			checkoutRepositoryMock: func(mc *minimock.Controller) domain.CheckoutRepository {
				mock := mocks.NewCheckoutRepositoryMock(mc)
				mock.AddToCartMock.Expect(ctx, user, sku, countSufficient).Return(nil)
				return mock
			},
			lomsClientMock: func(mc *minimock.Controller) domain.LOMSClient {
				mock := mocks.NewLOMSClientMock(mc)
				mock.StocksMock.Expect(ctx, sku).Return(nil, lomsErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			businessLogic := domain.New(tt.lomsClientMock(mc), nil, tt.checkoutRepositoryMock(mc), nil)
			err := businessLogic.AddToCart(tt.args.ctx, tt.args.user, tt.args.sku, tt.args.count)
			if tt.err != nil {
				require.ErrorContains(t, err, tt.err.Error())
			} else {
				require.Equal(t, tt.err, err)
			}
		})
	}
}

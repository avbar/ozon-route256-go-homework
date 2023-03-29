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

func TestDeleteFromCart(t *testing.T) {
	type args struct {
		ctx   context.Context
		user  int64
		sku   uint32
		count uint16
	}

	var (
		mc  = minimock.NewController(t)
		ctx = context.Background()

		user  = gofakeit.Int64()
		sku   = gofakeit.Uint32()
		count = gofakeit.Uint16()

		repoErr = errors.New("repo delete from cart error")
	)

	tests := []struct {
		name                   string
		args                   args
		err                    error
		checkoutRepositoryMock checkoutRepositoryMockFunc
	}{
		{
			name: "positive case",
			args: args{
				ctx:   ctx,
				user:  user,
				sku:   sku,
				count: count,
			},
			err: nil,
			checkoutRepositoryMock: func(mc *minimock.Controller) domain.CheckoutRepository {
				mock := mocks.NewCheckoutRepositoryMock(mc)
				mock.DeleteFromCartMock.Expect(ctx, user, sku, count).Return(nil)
				return mock
			},
		},

		{
			name: "repo error",
			args: args{
				ctx:   ctx,
				user:  user,
				sku:   sku,
				count: count,
			},
			err: repoErr,
			checkoutRepositoryMock: func(mc *minimock.Controller) domain.CheckoutRepository {
				mock := mocks.NewCheckoutRepositoryMock(mc)
				mock.DeleteFromCartMock.Expect(ctx, user, sku, count).Return(repoErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			businessLogic := domain.New(nil, nil, tt.checkoutRepositoryMock(mc), nil)
			err := businessLogic.DeleteFromCart(tt.args.ctx, tt.args.user, tt.args.sku, tt.args.count)
			if tt.err != nil {
				require.ErrorContains(t, err, tt.err.Error())
			} else {
				require.Equal(t, tt.err, err)
			}
		})
	}
}

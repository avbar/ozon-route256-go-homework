package domain_test

import (
	"context"
	"errors"
	"testing"

	"route256/checkout/internal/domain"
	"route256/checkout/internal/domain/mocks"
	"route256/libs/postgres/transactor"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
)

type lomsClientMockFunc func(mc *minimock.Controller) domain.LOMSClient
type productClientMockFunc func(mc *minimock.Controller) domain.ProductClient
type checkoutRepositoryMockFunc func(mc *minimock.Controller) domain.CheckoutRepository
type dbMockFunc func(mc *minimock.Controller) transactor.DB

func TestListCart(t *testing.T) {
	type args struct {
		ctx  context.Context
		user int64
	}

	var (
		mc  = minimock.NewController(t)
		ctx = context.Background()

		n     = 10
		user  = gofakeit.Int64()
		token = ""
		skus  = make([]uint32, 0, n)

		repoErr    = errors.New("read repo error")
		productErr = errors.New("product service error")

		repoRes    []domain.CartItem
		productRes = make(map[uint32]domain.Product)

		expectRes domain.Cart
	)

	for i := 0; i < n; i++ {
		sku := gofakeit.Uint32()
		skus = append(skus, sku)

		repoRes = append(repoRes, domain.CartItem{
			SKU:   sku,
			Count: gofakeit.Uint16(),
		})
	}

	for _, rr := range repoRes {
		productRes[rr.SKU] = domain.Product{
			Name:  gofakeit.BeerName(),
			Price: gofakeit.Uint32(),
		}

		expectRes.Items = append(expectRes.Items, domain.CartItemDetail{
			SKU:   rr.SKU,
			Count: rr.Count,
			Name:  productRes[rr.SKU].Name,
			Price: productRes[rr.SKU].Price,
		})
		expectRes.TotalPrice += productRes[rr.SKU].Price
	}

	tests := []struct {
		name                   string
		args                   args
		want                   domain.Cart
		err                    error
		checkoutRepositoryMock checkoutRepositoryMockFunc
		productClientMock      productClientMockFunc
	}{
		{
			name: "positive case",
			args: args{
				ctx:  ctx,
				user: user,
			},
			want: expectRes,
			err:  nil,
			checkoutRepositoryMock: func(mc *minimock.Controller) domain.CheckoutRepository {
				mock := mocks.NewCheckoutRepositoryMock(mc)
				mock.ListCartMock.Expect(ctx, user).Return(repoRes, nil)
				return mock
			},
			productClientMock: func(mc *minimock.Controller) domain.ProductClient {
				mock := mocks.NewProductClientMock(mc)
				mock.GetProductsMock.Expect(ctx, token, skus).Return(productRes, nil)
				return mock
			},
		},

		{
			name: "read repo error",
			args: args{
				ctx:  ctx,
				user: user,
			},
			want: domain.Cart{},
			err:  repoErr,
			checkoutRepositoryMock: func(mc *minimock.Controller) domain.CheckoutRepository {
				mock := mocks.NewCheckoutRepositoryMock(mc)
				mock.ListCartMock.Expect(ctx, user).Return(nil, repoErr)
				return mock
			},
			productClientMock: func(mc *minimock.Controller) domain.ProductClient {
				mock := mocks.NewProductClientMock(mc)
				mock.GetProductsMock.Expect(ctx, token, skus).Return(productRes, nil)
				return mock
			},
		},

		{
			name: "product service error",
			args: args{
				ctx:  ctx,
				user: user,
			},
			want: domain.Cart{},
			err:  productErr,
			checkoutRepositoryMock: func(mc *minimock.Controller) domain.CheckoutRepository {
				mock := mocks.NewCheckoutRepositoryMock(mc)
				mock.ListCartMock.Expect(ctx, user).Return(repoRes, nil)
				return mock
			},
			productClientMock: func(mc *minimock.Controller) domain.ProductClient {
				mock := mocks.NewProductClientMock(mc)
				mock.GetProductsMock.Expect(ctx, token, skus).Return(nil, productErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			businessLogic := domain.New(nil, tt.productClientMock(mc), tt.checkoutRepositoryMock(mc), nil)
			res, err := businessLogic.ListCart(ctx, user)
			require.Equal(t, tt.want, res)
			if tt.err != nil {
				require.ErrorContains(t, err, tt.err.Error())
			} else {
				require.Equal(t, tt.err, err)
			}
		})
	}
}

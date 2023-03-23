package domain_test

import (
	"context"
	"route256/checkout/internal/domain"
	"route256/checkout/internal/domain/mocks"
	"route256/libs/postgres/transactor"
	txMocks "route256/libs/postgres/transactor/mocks"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestPurchase(t *testing.T) {
	type args struct {
		ctx  context.Context
		user int64
	}

	var (
		mc  = minimock.NewController(t)
		ctx = context.Background()

		n    = 10
		user = gofakeit.Int64()

		opts = pgx.TxOptions{IsoLevel: pgx.RepeatableRead}
		tx   = txMocks.NewTxMock(t)

		listCartErr       = errors.New("list cart error")
		createOrderErr    = errors.New("create order error")
		deleteFromCartErr = errors.New("delete from cart error")
		commitErr         = errors.New("commit error")

		repoListCartRes    []domain.CartItem
		lomsCreateOrderReq []domain.OrderItem
		lomsCreateOrderRes = domain.OrderID(gofakeit.Int64())

		expectRes domain.OrderID = lomsCreateOrderRes
	)

	for i := 0; i < n; i++ {
		sku := gofakeit.Uint32()
		count := gofakeit.Uint16()

		repoListCartRes = append(repoListCartRes, domain.CartItem{
			SKU:   sku,
			Count: count,
		})

		lomsCreateOrderReq = append(lomsCreateOrderReq, domain.OrderItem{
			SKU:   sku,
			Count: count,
		})
	}

	tests := []struct {
		name                   string
		args                   args
		want                   domain.OrderID
		err                    error
		checkoutRepositoryMock checkoutRepositoryMockFunc
		lomsClientMock         lomsClientMockFunc
		dbMock                 dbMockFunc
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
				mock.ListCartMock.Expect(ctx, user).Return(repoListCartRes, nil)
				for _, r := range repoListCartRes {
					mock.DeleteFromCartMock.When(ctx, user, r.SKU, r.Count).Then(nil)
				}
				return mock
			},
			lomsClientMock: func(mc *minimock.Controller) domain.LOMSClient {
				mock := mocks.NewLOMSClientMock(mc)
				mock.CreateOrderMock.Expect(ctx, user, lomsCreateOrderReq).Return(lomsCreateOrderRes, nil)
				return mock
			},
			dbMock: func(mc *minimock.Controller) transactor.DB {
				mock := txMocks.NewDBMock(mc)
				mock.BeginTxMock.Expect(ctx, opts).Return(tx, nil)
				tx.CommitMock.Expect(ctx).Return(nil)
				return mock
			},
		},

		{
			name: "list cart error",
			args: args{
				ctx:  ctx,
				user: user,
			},
			want: 0,
			err:  listCartErr,
			checkoutRepositoryMock: func(mc *minimock.Controller) domain.CheckoutRepository {
				mock := mocks.NewCheckoutRepositoryMock(mc)
				mock.ListCartMock.Expect(ctx, user).Return(nil, listCartErr)
				for _, r := range repoListCartRes {
					mock.DeleteFromCartMock.When(ctx, user, r.SKU, r.Count).Then(nil)
				}
				return mock
			},
			lomsClientMock: func(mc *minimock.Controller) domain.LOMSClient {
				mock := mocks.NewLOMSClientMock(mc)
				mock.CreateOrderMock.Expect(ctx, user, lomsCreateOrderReq).Return(lomsCreateOrderRes, nil)
				return mock
			},
			dbMock: func(mc *minimock.Controller) transactor.DB {
				mock := txMocks.NewDBMock(mc)
				mock.BeginTxMock.Expect(ctx, opts).Return(tx, nil)
				tx.CommitMock.Expect(ctx).Return(nil)
				return mock
			},
		},

		{
			name: "create order error",
			args: args{
				ctx:  ctx,
				user: user,
			},
			want: 0,
			err:  createOrderErr,
			checkoutRepositoryMock: func(mc *minimock.Controller) domain.CheckoutRepository {
				mock := mocks.NewCheckoutRepositoryMock(mc)
				mock.ListCartMock.Expect(ctx, user).Return(repoListCartRes, nil)
				for _, r := range repoListCartRes {
					mock.DeleteFromCartMock.When(ctx, user, r.SKU, r.Count).Then(nil)
				}
				return mock
			},
			lomsClientMock: func(mc *minimock.Controller) domain.LOMSClient {
				mock := mocks.NewLOMSClientMock(mc)
				mock.CreateOrderMock.Expect(ctx, user, lomsCreateOrderReq).Return(0, createOrderErr)
				return mock
			},
			dbMock: func(mc *minimock.Controller) transactor.DB {
				mock := txMocks.NewDBMock(mc)
				mock.BeginTxMock.Expect(ctx, opts).Return(tx, nil)
				tx.RollbackMock.Expect(ctx).Return(nil)
				return mock
			},
		},

		{
			name: "delete from cart error",
			args: args{
				ctx:  ctx,
				user: user,
			},
			want: 0,
			err:  deleteFromCartErr,
			checkoutRepositoryMock: func(mc *minimock.Controller) domain.CheckoutRepository {
				mock := mocks.NewCheckoutRepositoryMock(mc)
				mock.ListCartMock.Expect(ctx, user).Return(repoListCartRes, nil)
				for _, r := range repoListCartRes {
					mock.DeleteFromCartMock.When(ctx, user, r.SKU, r.Count).Then(deleteFromCartErr)
				}
				return mock
			},
			lomsClientMock: func(mc *minimock.Controller) domain.LOMSClient {
				mock := mocks.NewLOMSClientMock(mc)
				mock.CreateOrderMock.Expect(ctx, user, lomsCreateOrderReq).Return(lomsCreateOrderRes, nil)
				return mock
			},
			dbMock: func(mc *minimock.Controller) transactor.DB {
				mock := txMocks.NewDBMock(mc)
				mock.BeginTxMock.Expect(ctx, opts).Return(tx, nil)
				tx.RollbackMock.Expect(ctx).Return(nil)
				return mock
			},
		},

		{
			name: "commit error",
			args: args{
				ctx:  ctx,
				user: user,
			},
			want: 0,
			err:  commitErr,
			checkoutRepositoryMock: func(mc *minimock.Controller) domain.CheckoutRepository {
				mock := mocks.NewCheckoutRepositoryMock(mc)
				mock.ListCartMock.Expect(ctx, user).Return(repoListCartRes, nil)
				for _, r := range repoListCartRes {
					mock.DeleteFromCartMock.When(ctx, user, r.SKU, r.Count).Then(nil)
				}
				return mock
			},
			lomsClientMock: func(mc *minimock.Controller) domain.LOMSClient {
				mock := mocks.NewLOMSClientMock(mc)
				mock.CreateOrderMock.Expect(ctx, user, lomsCreateOrderReq).Return(lomsCreateOrderRes, nil)
				return mock
			},
			dbMock: func(mc *minimock.Controller) transactor.DB {
				mock := txMocks.NewDBMock(mc)
				mock.BeginTxMock.Expect(ctx, opts).Return(tx, nil)
				tx.CommitMock.Expect(ctx).Return(commitErr)
				tx.RollbackMock.Expect(ctx).Return(nil)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			businessLogic := domain.New(tt.lomsClientMock(mc), nil, tt.checkoutRepositoryMock(mc), transactor.NewTransactionManager(tt.dbMock(mc)))
			res, err := businessLogic.Purchase(tt.args.ctx, tt.args.user)
			require.Equal(t, tt.want, res)
			if tt.err != nil {
				require.ErrorContains(t, err, tt.err.Error())
			} else {
				require.Equal(t, tt.err, err)
			}
		})
	}
}

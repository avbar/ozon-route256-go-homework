package domain_test

import (
	"context"
	"route256/libs/postgres/transactor"
	txMocks "route256/libs/postgres/transactor/mocks"
	"route256/loms/internal/domain"
	"route256/loms/internal/domain/mocks"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestCreateOrder(t *testing.T) {
	type args struct {
		ctx   context.Context
		user  int64
		items []domain.OrderItem
	}

	var (
		mc  = minimock.NewController(t)
		ctx = context.Background()

		opts = pgx.TxOptions{IsoLevel: pgx.RepeatableRead}
		tx   = txMocks.NewTxMock(t)

		n       = 10
		user    = gofakeit.Int64()
		items   []domain.OrderItem
		orderID = domain.OrderID(gofakeit.Int64())

		createOrderErr   = errors.New("repo create order error")
		saveOutboxErr    = errors.New("save to outbox error")
		createReserveErr = errors.New("repo create reserve error")
		changeStatusErr  = errors.New("repo change status error")
		commitErr        = errors.New("commit error")

		expectRes = orderID
	)

	for i := 0; i < n; i++ {
		items = append(items, domain.OrderItem{
			SKU:   gofakeit.Uint32(),
			Count: gofakeit.Uint16(),
		})
	}

	tests := []struct {
		name               string
		args               args
		want               domain.OrderID
		err                error
		lomsRepositoryMock lomsRepositoryMockFunc
		orderSenderMock    orderSenderMockFunc
		dbMock             dbMockFunc
	}{
		{
			name: "positive case",
			args: args{
				ctx:   ctx,
				user:  user,
				items: items,
			},
			want: expectRes,
			err:  nil,
			lomsRepositoryMock: func(mc *minimock.Controller) domain.LOMSRepository {
				mock := mocks.NewLOMSRepositoryMock(mc)
				mock.CreateOrderMock.Expect(ctx, user, items).Return(orderID, nil)
				mock.CreateReserveMock.Expect(ctx, orderID, items).Return(nil)
				mock.ChangeStatusMock.Expect(ctx, orderID, domain.OrderStatusAwaitingPayment).Return(nil)
				mock.SaveOrderToOutboxMock.When(ctx, orderID, domain.OrderStatusNew).Then(nil)
				mock.SaveOrderToOutboxMock.When(ctx, orderID, domain.OrderStatusAwaitingPayment).Then(nil)
				return mock
			},
			orderSenderMock: func(mc *minimock.Controller) domain.OrderSender {
				mock := mocks.NewOrderSenderMock(mc)
				mock.AddSuccessHandlerMock.Set(func(ctx context.Context, onSuccess func(orderID int64, status string)) {})
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
			name: "repo create order error",
			args: args{
				ctx:   ctx,
				user:  user,
				items: items,
			},
			want: 0,
			err:  createOrderErr,
			lomsRepositoryMock: func(mc *minimock.Controller) domain.LOMSRepository {
				mock := mocks.NewLOMSRepositoryMock(mc)
				mock.CreateOrderMock.Expect(ctx, user, items).Return(0, createOrderErr)
				mock.CreateReserveMock.Expect(ctx, orderID, items).Return(nil)
				mock.ChangeStatusMock.Expect(ctx, orderID, domain.OrderStatusAwaitingPayment).Return(nil)
				mock.SaveOrderToOutboxMock.When(ctx, orderID, domain.OrderStatusNew).Then(nil)
				mock.SaveOrderToOutboxMock.When(ctx, orderID, domain.OrderStatusAwaitingPayment).Then(nil)
				return mock
			},
			orderSenderMock: func(mc *minimock.Controller) domain.OrderSender {
				mock := mocks.NewOrderSenderMock(mc)
				mock.AddSuccessHandlerMock.Set(func(ctx context.Context, onSuccess func(orderID int64, status string)) {})
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
			name: "save to outbox error",
			args: args{
				ctx:   ctx,
				user:  user,
				items: items,
			},
			want: 0,
			err:  saveOutboxErr,
			lomsRepositoryMock: func(mc *minimock.Controller) domain.LOMSRepository {
				mock := mocks.NewLOMSRepositoryMock(mc)
				mock.CreateOrderMock.Expect(ctx, user, items).Return(orderID, nil)
				mock.CreateReserveMock.Expect(ctx, orderID, items).Return(nil)
				mock.ChangeStatusMock.Expect(ctx, orderID, domain.OrderStatusAwaitingPayment).Return(nil)
				mock.SaveOrderToOutboxMock.Expect(ctx, orderID, domain.OrderStatusNew).Return(saveOutboxErr)
				return mock
			},
			orderSenderMock: func(mc *minimock.Controller) domain.OrderSender {
				mock := mocks.NewOrderSenderMock(mc)
				mock.AddSuccessHandlerMock.Set(func(ctx context.Context, onSuccess func(orderID int64, status string)) {})
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
			name: "status failed after reserve error",
			args: args{
				ctx:   ctx,
				user:  user,
				items: items,
			},
			want: expectRes,
			err:  nil,
			lomsRepositoryMock: func(mc *minimock.Controller) domain.LOMSRepository {
				mock := mocks.NewLOMSRepositoryMock(mc)
				mock.CreateOrderMock.Expect(ctx, user, items).Return(orderID, nil)
				mock.CreateReserveMock.Expect(ctx, orderID, items).Return(createReserveErr)
				mock.ChangeStatusMock.Expect(ctx, orderID, domain.OrderStatusFailed).Return(nil)
				mock.SaveOrderToOutboxMock.When(ctx, orderID, domain.OrderStatusNew).Then(nil)
				mock.SaveOrderToOutboxMock.When(ctx, orderID, domain.OrderStatusFailed).Then(nil)
				return mock
			},
			orderSenderMock: func(mc *minimock.Controller) domain.OrderSender {
				mock := mocks.NewOrderSenderMock(mc)
				mock.AddSuccessHandlerMock.Set(func(ctx context.Context, onSuccess func(orderID int64, status string)) {})
				return mock
			},
			dbMock: func(mc *minimock.Controller) transactor.DB {
				mock := txMocks.NewDBMock(mc)
				mock.BeginTxMock.Expect(ctx, opts).Return(tx, nil)
				tx.CommitMock.Expect(ctx).Return(nil)
				tx.RollbackMock.Expect(ctx).Return(nil)
				return mock
			},
		},

		{
			name: "status failed after commit error",
			args: args{
				ctx:   ctx,
				user:  user,
				items: items,
			},
			want: expectRes,
			err:  nil,
			lomsRepositoryMock: func(mc *minimock.Controller) domain.LOMSRepository {
				mock := mocks.NewLOMSRepositoryMock(mc)
				mock.CreateOrderMock.Expect(ctx, user, items).Return(orderID, nil)
				mock.CreateReserveMock.Expect(ctx, orderID, items).Return(nil)
				mock.ChangeStatusMock.Expect(ctx, orderID, domain.OrderStatusFailed).Return(nil)
				mock.SaveOrderToOutboxMock.When(ctx, orderID, domain.OrderStatusNew).Then(nil)
				mock.SaveOrderToOutboxMock.When(ctx, orderID, domain.OrderStatusFailed).Then(nil)
				return mock
			},
			orderSenderMock: func(mc *minimock.Controller) domain.OrderSender {
				mock := mocks.NewOrderSenderMock(mc)
				mock.AddSuccessHandlerMock.Set(func(ctx context.Context, onSuccess func(orderID int64, status string)) {})
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

		{
			name: "change status error",
			args: args{
				ctx:   ctx,
				user:  user,
				items: items,
			},
			want: expectRes,
			err:  changeStatusErr,
			lomsRepositoryMock: func(mc *minimock.Controller) domain.LOMSRepository {
				mock := mocks.NewLOMSRepositoryMock(mc)
				mock.CreateOrderMock.Expect(ctx, user, items).Return(orderID, nil)
				mock.CreateReserveMock.Expect(ctx, orderID, items).Return(nil)
				mock.ChangeStatusMock.Expect(ctx, orderID, domain.OrderStatusAwaitingPayment).Return(changeStatusErr)
				mock.SaveOrderToOutboxMock.When(ctx, orderID, domain.OrderStatusNew).Then(nil)
				mock.SaveOrderToOutboxMock.When(ctx, orderID, domain.OrderStatusAwaitingPayment).Then(nil)
				return mock
			},
			orderSenderMock: func(mc *minimock.Controller) domain.OrderSender {
				mock := mocks.NewOrderSenderMock(mc)
				mock.AddSuccessHandlerMock.Set(func(ctx context.Context, onSuccess func(orderID int64, status string)) {})
				return mock
			},
			dbMock: func(mc *minimock.Controller) transactor.DB {
				mock := txMocks.NewDBMock(mc)
				mock.BeginTxMock.Expect(ctx, opts).Return(tx, nil)
				tx.CommitMock.Expect(ctx).Return(nil)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			businessLogic := domain.New(tt.lomsRepositoryMock(mc), transactor.NewTransactionManager(tt.dbMock(mc)), tt.orderSenderMock(mc))
			res, err := businessLogic.CreateOrder(tt.args.ctx, tt.args.user, tt.args.items)
			require.Equal(t, tt.want, res)
			if tt.err != nil {
				require.ErrorContains(t, err, tt.err.Error())
			} else {
				require.Equal(t, tt.err, err)
			}
		})
	}
}

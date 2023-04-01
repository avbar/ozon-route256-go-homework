package domain_test

import (
	"context"
	"errors"
	"route256/libs/postgres/transactor"
	txMocks "route256/libs/postgres/transactor/mocks"
	"route256/loms/internal/domain"
	"route256/loms/internal/domain/mocks"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/require"
)

func TestCancelOrder(t *testing.T) {
	type args struct {
		ctx     context.Context
		orderID domain.OrderID
	}

	var (
		mc  = minimock.NewController(t)
		ctx = context.Background()

		opts = pgx.TxOptions{IsoLevel: pgx.RepeatableRead}
		tx   = txMocks.NewTxMock(t)

		orderID = domain.OrderID(gofakeit.Int64())

		getStatusErr     = errors.New("repo get status error")
		wrongStatusErr   = domain.ErrWrongStatus
		changeStatusErr  = errors.New("repo change status error")
		cancelReserveErr = errors.New("repo cancel reserve error")
		commitErr        = errors.New("commit error")
	)

	tests := []struct {
		name               string
		args               args
		err                error
		lomsRepositoryMock lomsRepositoryMockFunc
		orderSenderMock    orderSenderMockFunc
		dbMock             dbMockFunc
	}{
		{
			name: "positive case",
			args: args{
				ctx:     ctx,
				orderID: orderID,
			},
			err: nil,
			lomsRepositoryMock: func(mc *minimock.Controller) domain.LOMSRepository {
				mock := mocks.NewLOMSRepositoryMock(mc)
				mock.GetStatusMock.Expect(ctx, orderID).Return(domain.OrderStatusAwaitingPayment, nil)
				mock.CancelReserveMock.Expect(ctx, orderID).Return(nil)
				mock.ChangeStatusMock.Expect(ctx, orderID, domain.OrderStatusCancelled).Return(nil)
				mock.SaveOrderToOutboxMock.Expect(ctx, orderID, domain.OrderStatusCancelled).Return(nil)
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
			name: "get status error",
			args: args{
				ctx:     ctx,
				orderID: orderID,
			},
			err: getStatusErr,
			lomsRepositoryMock: func(mc *minimock.Controller) domain.LOMSRepository {
				mock := mocks.NewLOMSRepositoryMock(mc)
				mock.GetStatusMock.Expect(ctx, orderID).Return("", getStatusErr)
				mock.CancelReserveMock.Expect(ctx, orderID).Return(nil)
				mock.ChangeStatusMock.Expect(ctx, orderID, domain.OrderStatusCancelled).Return(nil)
				mock.SaveOrderToOutboxMock.Expect(ctx, orderID, domain.OrderStatusCancelled).Return(nil)
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
			name: "wrong status",
			args: args{
				ctx:     ctx,
				orderID: orderID,
			},
			err: wrongStatusErr,
			lomsRepositoryMock: func(mc *minimock.Controller) domain.LOMSRepository {
				mock := mocks.NewLOMSRepositoryMock(mc)
				mock.GetStatusMock.Expect(ctx, orderID).Return(domain.OrderStatusPayed, nil)
				mock.CancelReserveMock.Expect(ctx, orderID).Return(nil)
				mock.ChangeStatusMock.Expect(ctx, orderID, domain.OrderStatusCancelled).Return(nil)
				mock.SaveOrderToOutboxMock.Expect(ctx, orderID, domain.OrderStatusCancelled).Return(nil)
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
			name: "change status error",
			args: args{
				ctx:     ctx,
				orderID: orderID,
			},
			err: changeStatusErr,
			lomsRepositoryMock: func(mc *minimock.Controller) domain.LOMSRepository {
				mock := mocks.NewLOMSRepositoryMock(mc)
				mock.GetStatusMock.Expect(ctx, orderID).Return(domain.OrderStatusAwaitingPayment, nil)
				mock.CancelReserveMock.Expect(ctx, orderID).Return(nil)
				mock.ChangeStatusMock.Expect(ctx, orderID, domain.OrderStatusCancelled).Return(changeStatusErr)
				mock.SaveOrderToOutboxMock.Expect(ctx, orderID, domain.OrderStatusCancelled).Return(nil)
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
			name: "cancel reserve error",
			args: args{
				ctx:     ctx,
				orderID: orderID,
			},
			err: cancelReserveErr,
			lomsRepositoryMock: func(mc *minimock.Controller) domain.LOMSRepository {
				mock := mocks.NewLOMSRepositoryMock(mc)
				mock.GetStatusMock.Expect(ctx, orderID).Return(domain.OrderStatusAwaitingPayment, nil)
				mock.CancelReserveMock.Expect(ctx, orderID).Return(cancelReserveErr)
				mock.ChangeStatusMock.Expect(ctx, orderID, domain.OrderStatusCancelled).Return(nil)
				mock.SaveOrderToOutboxMock.Expect(ctx, orderID, domain.OrderStatusCancelled).Return(nil)
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
			name: "commit error",
			args: args{
				ctx:     ctx,
				orderID: orderID,
			},
			err: commitErr,
			lomsRepositoryMock: func(mc *minimock.Controller) domain.LOMSRepository {
				mock := mocks.NewLOMSRepositoryMock(mc)
				mock.GetStatusMock.Expect(ctx, orderID).Return(domain.OrderStatusAwaitingPayment, nil)
				mock.CancelReserveMock.Expect(ctx, orderID).Return(nil)
				mock.ChangeStatusMock.Expect(ctx, orderID, domain.OrderStatusCancelled).Return(nil)
				mock.SaveOrderToOutboxMock.Expect(ctx, orderID, domain.OrderStatusCancelled).Return(nil)
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			businessLogic := domain.New(tt.lomsRepositoryMock(mc), transactor.NewTransactionManager(tt.dbMock(mc)), tt.orderSenderMock(mc))
			err := businessLogic.CancelOrder(tt.args.ctx, tt.args.orderID)
			if tt.err != nil {
				require.ErrorContains(t, err, tt.err.Error())
			} else {
				require.Equal(t, tt.err, err)
			}
		})
	}
}

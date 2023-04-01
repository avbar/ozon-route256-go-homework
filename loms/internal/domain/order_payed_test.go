package domain_test

import (
	"context"
	"route256/loms/internal/domain"
	"route256/loms/internal/domain/mocks"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestOrderPayed(t *testing.T) {
	type args struct {
		ctx     context.Context
		orderID domain.OrderID
	}

	var (
		mc  = minimock.NewController(t)
		ctx = context.Background()

		orderID = domain.OrderID(gofakeit.Int64())

		deleteReserveErr = errors.New("delete reserve error")
		changeStatusErr  = errors.New("change status error")
	)

	tests := []struct {
		name               string
		args               args
		err                error
		lomsRepositoryMock lomsRepositoryMockFunc
		orderSenderMock    orderSenderMockFunc
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
				mock.DeleteReserveMock.Expect(ctx, orderID).Return(nil)
				mock.ChangeStatusMock.Expect(ctx, orderID, domain.OrderStatusPayed).Return(nil)
				mock.SaveOrderToOutboxMock.Expect(ctx, orderID, domain.OrderStatusPayed).Return(nil)
				return mock
			},
			orderSenderMock: func(mc *minimock.Controller) domain.OrderSender {
				mock := mocks.NewOrderSenderMock(mc)
				mock.AddSuccessHandlerMock.Set(func(ctx context.Context, onSuccess func(orderID int64, status string)) {})
				return mock
			},
		},

		{
			name: "delete reserve error",
			args: args{
				ctx:     ctx,
				orderID: orderID,
			},
			err: deleteReserveErr,
			lomsRepositoryMock: func(mc *minimock.Controller) domain.LOMSRepository {
				mock := mocks.NewLOMSRepositoryMock(mc)
				mock.DeleteReserveMock.Expect(ctx, orderID).Return(deleteReserveErr)
				mock.ChangeStatusMock.Expect(ctx, orderID, domain.OrderStatusPayed).Return(nil)
				mock.SaveOrderToOutboxMock.Expect(ctx, orderID, domain.OrderStatusPayed).Return(nil)
				return mock
			},
			orderSenderMock: func(mc *minimock.Controller) domain.OrderSender {
				mock := mocks.NewOrderSenderMock(mc)
				mock.AddSuccessHandlerMock.Set(func(ctx context.Context, onSuccess func(orderID int64, status string)) {})
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
				mock.DeleteReserveMock.Expect(ctx, orderID).Return(nil)
				mock.ChangeStatusMock.Expect(ctx, orderID, domain.OrderStatusPayed).Return(changeStatusErr)
				mock.SaveOrderToOutboxMock.Expect(ctx, orderID, domain.OrderStatusPayed).Return(nil)
				return mock
			},
			orderSenderMock: func(mc *minimock.Controller) domain.OrderSender {
				mock := mocks.NewOrderSenderMock(mc)
				mock.AddSuccessHandlerMock.Set(func(ctx context.Context, onSuccess func(orderID int64, status string)) {})
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			businessLogic := domain.New(tt.lomsRepositoryMock(mc), nil, tt.orderSenderMock(mc))
			err := businessLogic.OrderPayed(tt.args.ctx, tt.args.orderID)
			if tt.err != nil {
				require.ErrorContains(t, err, tt.err.Error())
			} else {
				require.Equal(t, tt.err, err)
			}
		})
	}
}

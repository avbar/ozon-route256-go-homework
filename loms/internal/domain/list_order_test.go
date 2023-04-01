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

func TestListOrder(t *testing.T) {
	type args struct {
		ctx     context.Context
		orderID domain.OrderID
	}

	var (
		mc  = minimock.NewController(t)
		ctx = context.Background()

		n       = 10
		orderID = domain.OrderID(gofakeit.Int64())

		repoErr = errors.New("read repo error")

		repoRes   domain.Order
		expectRes domain.Order
	)

	repoRes.Status = domain.OrderStatusAwaitingPayment
	repoRes.User = gofakeit.Int64()
	for i := 0; i < n; i++ {
		repoRes.Items = append(repoRes.Items, domain.OrderItem{
			SKU:   gofakeit.Uint32(),
			Count: gofakeit.Uint16(),
		})
	}
	expectRes = repoRes

	tests := []struct {
		name               string
		args               args
		want               domain.Order
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
			want: expectRes,
			err:  nil,
			lomsRepositoryMock: func(mc *minimock.Controller) domain.LOMSRepository {
				mock := mocks.NewLOMSRepositoryMock(mc)
				mock.ListOrderMock.Expect(ctx, orderID).Return(repoRes, nil)
				return mock
			},
			orderSenderMock: func(mc *minimock.Controller) domain.OrderSender {
				mock := mocks.NewOrderSenderMock(mc)
				mock.AddSuccessHandlerMock.Set(func(ctx context.Context, onSuccess func(orderID int64, status string)) {})
				return mock
			},
		},

		{
			name: "read repo error",
			args: args{
				ctx:     ctx,
				orderID: orderID,
			},
			want: domain.Order{},
			err:  repoErr,
			lomsRepositoryMock: func(mc *minimock.Controller) domain.LOMSRepository {
				mock := mocks.NewLOMSRepositoryMock(mc)
				mock.ListOrderMock.Expect(ctx, orderID).Return(domain.Order{}, repoErr)
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
			res, err := businessLogic.ListOrder(tt.args.ctx, tt.args.orderID)
			require.Equal(t, tt.want, res)
			if tt.err != nil {
				require.ErrorContains(t, err, tt.err.Error())
			} else {
				require.Equal(t, tt.err, err)
			}
		})
	}
}

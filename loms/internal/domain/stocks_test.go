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

func TestStocks(t *testing.T) {
	type args struct {
		ctx context.Context
		sku uint32
	}

	var (
		mc  = minimock.NewController(t)
		ctx = context.Background()

		n   = 10
		sku = gofakeit.Uint32()

		repoErr = errors.New("read repo error")

		repoRes   []domain.Stock
		expectRes []domain.Stock
	)

	for i := 0; i < n; i++ {
		repoRes = append(repoRes, domain.Stock{
			WarehouseID: gofakeit.Int64(),
			Count:       gofakeit.Uint64(),
		})
	}
	expectRes = repoRes

	tests := []struct {
		name               string
		args               args
		want               []domain.Stock
		err                error
		lomsRepositoryMock lomsRepositoryMockFunc
	}{
		{
			name: "positive case",
			args: args{
				ctx: ctx,
				sku: sku,
			},
			want: expectRes,
			err:  nil,
			lomsRepositoryMock: func(mc *minimock.Controller) domain.LOMSRepository {
				mock := mocks.NewLOMSRepositoryMock(mc)
				mock.StocksMock.Expect(ctx, sku).Return(repoRes, nil)
				return mock
			},
		},

		{
			name: "read repo error",
			args: args{
				ctx: ctx,
				sku: sku,
			},
			want: nil,
			err:  repoErr,
			lomsRepositoryMock: func(mc *minimock.Controller) domain.LOMSRepository {
				mock := mocks.NewLOMSRepositoryMock(mc)
				mock.StocksMock.Expect(ctx, sku).Return(nil, repoErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			businessLogic := domain.New(tt.lomsRepositoryMock(mc), nil)
			res, err := businessLogic.Stocks(tt.args.ctx, tt.args.sku)
			require.Equal(t, tt.want, res)
			if tt.err != nil {
				require.ErrorContains(t, err, tt.err.Error())
			} else {
				require.Equal(t, tt.err, err)
			}
		})
	}
}

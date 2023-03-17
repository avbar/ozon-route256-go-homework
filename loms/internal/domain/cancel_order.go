package domain

import (
	"context"

	"github.com/pkg/errors"
)

var (
	ErrWrongStatus = errors.New("wrong status")
)

func (m *Model) CancelOrder(ctx context.Context, orderID OrderID) error {
	status, err := m.lomsRepository.GetStatus(ctx, orderID)
	if err != nil {
		return err
	}

	if (status != OrderStatusNew) && (status != OrderStatusAwaitingPayment) {
		return errors.WithMessage(ErrWrongStatus, "cancelling order")
	}

	err = m.transactionManager.RunRepeatableRead(ctx, func(ctxTX context.Context) error {
		return m.lomsRepository.CancelReserve(ctx, orderID)
	})
	if err != nil {
		return err
	}

	return m.lomsRepository.ChangeStatus(ctx, orderID, OrderStatusCancelled)
}

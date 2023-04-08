package domain

import "context"

func (m *Model) ChangeStatus(ctx context.Context, orderID OrderID, status string) error {
	err := m.lomsRepository.ChangeStatus(ctx, orderID, status)
	if err != nil {
		return err
	}
	return m.lomsRepository.SaveOrderToOutbox(ctx, orderID, status)
}

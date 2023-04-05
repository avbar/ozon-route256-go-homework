package domain

import (
	"context"
	"route256/libs/logger"
	"route256/libs/workerpool"
	"sync"
	"time"

	"go.uber.org/zap"
)

const (
	orderTimeout  time.Duration = time.Minute * 10
	amountWorkers               = 5
)

// Создаёт Ticker для проверки заказов 1 раз в минуту
func (m *Model) runOrderTimer(ctx context.Context) {
	ticker := time.NewTicker(time.Minute * 1)

	for {
		select {
		case <-ticker.C:
			m.cancelOldOrders(ctx)
		case <-ctx.Done():
			return
		}
	}
}

// Отменяет заказы, не оплаченные в течение времени orderTimeout
func (m *Model) cancelOldOrders(ctx context.Context) {
	orderTime := time.Now().Add(-orderTimeout)
	logger.Info("cancelling unpaid orders created before time", zap.Time("time", orderTime))

	orders, err := m.lomsRepository.GetOldOrders(ctx, orderTime)
	if err != nil {
		logger.Error("error getting old orders", zap.Error(err))
		return
	}
	if len(orders) == 0 {
		return
	}

	// Создаём WorkerPool размера amountWorkers
	workerPool := workerpool.NewPool[error](ctx, amountWorkers)

	var wgSubmit sync.WaitGroup
	for _, orderID := range orders {
		wgSubmit.Add(1)

		// Добавляем в WorkerPool запрос на отмену заказа
		go func(orderID OrderID) {
			defer wgSubmit.Done()

			workerPool.Submit(ctx, workerpool.Task[error]{
				Callback: func() error {
					err := m.CancelOrder(ctx, orderID)
					if err != nil {
						logger.Error("error cancelling old order", zap.Int64("order id", int64(orderID)), zap.Error(err))
					} else {
						logger.Info("order cancelled", zap.Int64("order id", int64(orderID)))
					}
					return err
				},
			})
		}(orderID)
	}

	go workerPool.SkipOutput(ctx)

	// Дожидаемся окончания работы запросов на отмену
	wgSubmit.Wait()
	workerPool.Close()
}

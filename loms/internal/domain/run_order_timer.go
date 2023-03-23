package domain

import (
	"context"
	"log"
	"route256/libs/workerpool"
	"sync"
	"time"
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
	log.Printf("cancelling unpaid orders created before %v", orderTime.Format(time.DateTime))

	orders, err := m.lomsRepository.GetOldOrders(ctx, orderTime)
	if err != nil {
		log.Printf("error getting old orders: %v", err)
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
						log.Printf("error cancelling old order %v: %v", orderID, err)
					} else {
						log.Printf("order %v cancelled", orderID)
					}
					return err
				},
			})
		}(orderID)
	}

	// Дожидаемся окончания работы запросов на отмену
	wgSubmit.Wait()
	workerPool.Close()
}

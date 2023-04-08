package domain

import (
	"context"
	"log"
	"route256/libs/workerpool"
	"sync"
	"time"
)

const (
	outboxWorkers = 5
)

// Создаёт Ticker для отправки из Outbox статусов заказов
func (m *Model) runOutboxProcessor(ctx context.Context) {
	ticker := time.NewTicker(time.Second * 10)

	for {
		select {
		case <-ticker.C:
			m.sendOrderStatuses(ctx)
		case <-ctx.Done():
			return
		}
	}
}

// Отправляет статусы заказов в Kafka
func (m *Model) sendOrderStatuses(ctx context.Context) {
	orders, err := m.lomsRepository.GetOrdersFromOutbox(ctx)
	if err != nil {
		log.Printf("error reading outbox: %v", err)
		return
	}
	if len(orders) == 0 {
		return
	}

	// Создаём WorkerPool размера outboxWorkers
	workerPool := workerpool.NewPool[error](ctx, outboxWorkers)

	var wgSubmit sync.WaitGroup
	for _, order := range orders {
		wgSubmit.Add(1)

		// Добавляем в WorkerPool статус заказа
		go func(orderID OrderID, status string) {
			defer wgSubmit.Done()

			workerPool.Submit(ctx, workerpool.Task[error]{
				Callback: func() error {
					m.orderSender.SendOrderStatus(ctx, int64(orderID), status)
					return nil
				},
			})
		}(order.OrderID, order.Status)
	}

	// Дожидаемся окончания работы горутин
	wgSubmit.Wait()
	workerPool.Close()
}

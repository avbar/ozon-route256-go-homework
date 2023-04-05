package workerpool

import (
	"context"
	"sync"
)

// Задачи для обработки
type Task[Out any] struct {
	Callback func() Out
}

// WorkerPool
type WorkerPool[Out any] interface {
	Submit(context.Context, Task[Out])
	OutQueue() <-chan Out
	SkipOutput(context.Context)
	Close()
}

var _ WorkerPool[any] = &pool[any]{}

type pool[Out any] struct {
	// число обработчиков
	amountWorkers int

	// WaitGroup для синхронизации обработчиков
	wg sync.WaitGroup

	// входной канал с задачами
	taskQueue chan Task[Out]
	// выходной канал с результатами работы
	outQueue chan Out
}

// Создание пула
func NewPool[Out any](ctx context.Context, amountWorkers int) *pool[Out] {
	p := &pool[Out]{
		amountWorkers: amountWorkers,
	}

	p.init(ctx)

	return p
}

// Добавление задач во входной канал
func (p *pool[Out]) Submit(ctx context.Context, task Task[Out]) {
	select {
	case <-ctx.Done():
		return
	case p.taskQueue <- task:
	}
}

// Получение выходного канала
func (p *pool[Out]) OutQueue() <-chan Out {
	return p.outQueue
}

// Очистка выходного канала, если результат не нужен
func (p *pool[Out]) SkipOutput(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-p.outQueue:
		}
	}
}

// Закрытие каналов
func (p *pool[Out]) Close() {
	// Закрываем входной канал
	close(p.taskQueue)

	// Дожидаемся окончания работы всех обработчиков
	p.wg.Wait()

	// Закрываем выходной канал
	close(p.outQueue)
}

// Инициализация каналов
func (p *pool[Out]) init(ctx context.Context) {
	// Создаём входной и выходной каналы с буфером размера amountWorkers
	p.taskQueue = make(chan Task[Out], p.amountWorkers)
	p.outQueue = make(chan Out, p.amountWorkers)

	// Запускаем amountWorkers штук обработчиков
	for i := 0; i < p.amountWorkers; i++ {
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			worker(ctx, p.taskQueue, p.outQueue)
		}()
	}
}

// Берём задачи из канала taskQueue, выполняем, результат помещаем в канал outQueue
func worker[Out any](ctx context.Context, taskQueue <-chan Task[Out], outQueue chan<- Out) {
	for task := range taskQueue {
		select {
		case <-ctx.Done():
			return
		case outQueue <- task.Callback():
		}
	}
}

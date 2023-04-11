package cache

import (
	"container/list"
	"context"
	"sync"
	"time"
)

// Структура данных, хранящихся в кэше
type cacheElement struct {
	// ключ
	key string
	// значение
	value interface{}
	// срок действия
	validUntil time.Time
}

// Кэш с поддержкой TTL и LRU
type Cache struct {
	// данные кэша
	// в мапе хранятся указатели на элементы списка lru
	// структура cacheElement сохраняется в list.Element.Value
	data map[string]*list.Element
	// список для реализации LRU
	lru *list.List
	// Mutex для data/lru
	mx sync.Mutex
	// TTL
	ttl time.Duration
	// при достижении размера capacity последний элемент будет вытесняться из кэша
	capacity int
}

func New(ttl time.Duration, capacity int) *Cache {
	c := &Cache{
		data:     make(map[string]*list.Element),
		lru:      list.New(),
		ttl:      ttl,
		capacity: capacity,
	}

	go c.runCacheClearing(context.Background())

	return c
}

func (c *Cache) Set(ctx context.Context, key string, value interface{}) {
	// новое значение в кэше
	cacheElem := cacheElement{
		key:        key,
		value:      value,
		validUntil: time.Now().Add(c.ttl),
	}

	c.mx.Lock()
	defer c.mx.Unlock()

	if listElem, ok := c.data[key]; ok {
		// ключ уже есть в кэше
		// обновляем значение
		listElem.Value = cacheElem
		// поднимаем запись на верх списка
		c.lru.MoveToFront(listElem)
		return
	}

	// новый ключ
	if c.lru.Len() == c.capacity {
		// кэш заполнен
		// удаляем последнюю запись в списке
		lastElem := c.lru.Back()
		delete(c.data, lastElem.Value.(cacheElement).key)
		c.lru.Remove(lastElem)
	}
	// добавляем новый элемент
	c.data[key] = c.lru.PushFront(cacheElem)
}

func (c *Cache) Get(ctx context.Context, key string) (interface{}, bool) {
	var res interface{}

	timeStart := time.Now()
	// получаем данные по ключу
	c.mx.Lock()
	listElem, ok := c.data[key]
	if ok {
		// если ключ есть в кэше, то поднимаем запись наверх
		c.lru.MoveToFront(listElem)
		res = listElem.Value.(cacheElement).value
	}
	c.mx.Unlock()
	elapsedTime := time.Since(timeStart)

	// метрики
	HistogramResponseTime.Observe(elapsedTime.Seconds())

	if ok {
		HitsCounter.Inc()
	} else {
		ErrorsCounter.Inc()
	}

	return res, ok
}

// Создаёт Ticker для запуска очистки кэша по TTL
func (c *Cache) runCacheClearing(ctx context.Context) {
	ticker := time.NewTicker(time.Second * 10)

	for {
		select {
		case <-ticker.C:
			c.CheckExpiration(ctx)
		case <-ctx.Done():
			return
		}
	}
}

// Удаляет из кэша данные с истекшим сроком
func (c *Cache) CheckExpiration(ctx context.Context) {
	currTime := time.Now()

	c.mx.Lock()
	defer c.mx.Unlock()

	for key, listElem := range c.data {
		if listElem.Value.(cacheElement).validUntil.Before(currTime) {
			delete(c.data, key)
			c.lru.Remove(listElem)
		}
	}
}

package productservice

import (
	"context"
	"fmt"
	"route256/checkout/internal/domain"
	productAPI "route256/checkout/pkg/product_service_v1"
	"route256/libs/cache"
	"route256/libs/workerpool"
	"sync"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
)

var _ Client = (*client)(nil)

type ProductResult struct {
	SKU     uint32
	Product domain.Product
	Err     error
}

type Client interface {
	GetProduct(ctx context.Context, token string, sku uint32) (domain.Product, error)
	GetProducts(ctx context.Context, token string, skus []uint32) (map[uint32]domain.Product, error)
}

type client struct {
	productClient productAPI.ProductServiceClient
	limiter       rate.Limiter
	cache         cache.Cache
}

func NewClient(cc *grpc.ClientConn) *client {
	c := &client{
		productClient: productAPI.NewProductServiceClient(cc),
	}

	c.limiter = *rate.NewLimiter(rate.Every(time.Second/10), 10)

	c.cache = *cache.New(time.Minute*10, 50)

	return c
}

func (c *client) getProductFromService(ctx context.Context, token string, sku uint32) (domain.Product, error) {
	c.limiter.Wait(ctx)

	res, err := c.productClient.GetProduct(ctx, &productAPI.GetProductRequest{
		Token: token,
		Sku:   sku,
	})
	if err != nil {
		return domain.Product{}, errors.Wrap(err, "calling GetProduct")
	}

	product := domain.Product{
		Name:  res.GetName(),
		Price: res.GetPrice(),
	}

	return product, nil
}

func (c *client) getProductFromCache(ctx context.Context, sku uint32) (domain.Product, bool) {
	var product domain.Product

	value, ok := c.cache.Get(ctx, fmt.Sprint(sku))
	if ok {
		product = value.(domain.Product)
	}
	return product, ok
}

func (c *client) cacheProduct(ctx context.Context, sku uint32, product domain.Product) {
	c.cache.Set(ctx, fmt.Sprint(sku), product)
}

func (c *client) GetProduct(ctx context.Context, token string, sku uint32) (domain.Product, error) {
	var err error
	product, ok := c.getProductFromCache(ctx, sku)
	if !ok {
		product, err = c.getProductFromService(ctx, token, sku)
		if err != nil {
			return domain.Product{}, err
		}

		c.cacheProduct(ctx, sku, product)
	}

	return product, nil
}

func (c *client) GetProducts(ctx context.Context, token string, skus []uint32) (map[uint32]domain.Product, error) {
	const amountWorkers = 5

	// Создаём WorkerPool размера amountWorkers
	workerPool := workerpool.NewPool[ProductResult](ctx, amountWorkers)

	var wgSubmit sync.WaitGroup
	for _, sku := range skus {
		wgSubmit.Add(1)

		// Добавляем в WorkerPool запрос на получение данных по продукту
		go func(sku uint32) {
			defer wgSubmit.Done()

			workerPool.Submit(ctx, workerpool.Task[ProductResult]{
				Callback: func() ProductResult {
					product, err := c.GetProduct(ctx, token, sku)
					if err != nil {
						return ProductResult{
							SKU:     sku,
							Product: domain.Product{},
							Err:     err,
						}
					}
					return ProductResult{
						SKU:     sku,
						Product: product,
						Err:     nil,
					}
				},
			})
		}(sku)
	}

	var productResult []ProductResult = make([]ProductResult, 0)

	var wgResult sync.WaitGroup
	wgResult.Add(1)
	// Собираем полученные описания продуктов
	go func() {
		defer wgResult.Done()

		for product := range workerPool.OutQueue() {
			productResult = append(productResult, product)
		}
	}()

	// Дожидаемся окончания работы с ProductService
	wgSubmit.Wait()
	workerPool.Close()
	wgResult.Wait()

	products := make(map[uint32]domain.Product)

	// Формируем отображение с описаниями продуктов
	for _, p := range productResult {
		if p.Err != nil {
			return nil, p.Err
		}
		products[p.SKU] = p.Product
	}

	return products, nil
}

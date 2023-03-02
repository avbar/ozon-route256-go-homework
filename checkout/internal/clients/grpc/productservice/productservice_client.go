package productservice

import (
	"context"
	"route256/checkout/internal/domain"
	productAPI "route256/checkout/pkg/product_service_v1"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

var _ Client = (*client)(nil)

type Client interface {
	GetProduct(ctx context.Context, token string, sku uint32) (domain.Product, error)
}

type client struct {
	productClient productAPI.ProductServiceClient
}

func NewClient(cc *grpc.ClientConn) *client {
	return &client{
		productClient: productAPI.NewProductServiceClient(cc),
	}
}

func (c *client) GetProduct(ctx context.Context, token string, sku uint32) (domain.Product, error) {
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

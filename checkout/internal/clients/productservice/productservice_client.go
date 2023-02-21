package productservice

import (
	"context"
	"route256/checkout/internal/domain"
	"route256/libs/clntwrapper"

	"github.com/pkg/errors"
)

type Client struct {
	url string

	urlGetProduct string
}

func New(url string) *Client {
	return &Client{
		url: url,

		urlGetProduct: url + "/get_product",
	}
}

type ProductRequest struct {
	Token string `json:"token"`
	SKU   uint32 `json:"sku"`
}

type ProductResponse struct {
	Name  string `json:"name"`
	Price uint32 `json:"price"`
}

func (c *Client) GetProduct(ctx context.Context, token string, sku uint32) (domain.Product, error) {
	request := ProductRequest{Token: token, SKU: sku}

	var product domain.Product

	wrapper := clntwrapper.New[ProductRequest, ProductResponse](c.urlGetProduct)
	response, err := wrapper.Handle(ctx, request)
	if err != nil {
		return product, errors.Wrap(err, "calling client")
	}

	product = domain.Product{
		Name:  response.Name,
		Price: response.Price,
	}

	return product, nil
}

package clntwrapper

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

type Client[Req any, Res any] struct {
	url string
}

func New[Req any, Res any](url string) *Client[Req, Res] {
	return &Client[Req, Res]{
		url: url,
	}
}

func (c *Client[Req, Res]) Handle(ctx context.Context, request Req) (*Res, error) {
	rawJSON, err := json.Marshal(request)
	if err != nil {
		return nil, errors.Wrap(err, "marshaling json")
	}

	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url, bytes.NewBuffer(rawJSON))
	if err != nil {
		return nil, errors.Wrap(err, "creating http request")
	}

	httpResponse, err := http.DefaultClient.Do(httpRequest)
	if err != nil {
		return nil, errors.Wrap(err, "calling http")
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("wrong status code: %d", httpResponse.StatusCode)
	}

	var response Res
	err = json.NewDecoder(httpResponse.Body).Decode(&response)
	if err != nil {
		return nil, errors.Wrap(err, "decoding json")
	}

	return &response, nil
}

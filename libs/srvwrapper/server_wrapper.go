package srvwrapper

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

type Validator interface {
	Validate() error
}

type Wrapper[Req Validator, Res any] struct {
	fn func(ctx context.Context, req Req) (Res, error)
}

func New[Req Validator, Res any](fn func(ctx context.Context, req Req) (Res, error)) *Wrapper[Req, Res] {
	return &Wrapper[Req, Res]{
		fn: fn,
	}
}

func (w *Wrapper[Req, Res]) ServeHTTP(resWriter http.ResponseWriter, httpReq *http.Request) {
	ctx := httpReq.Context()

	limitedReader := io.LimitReader(httpReq.Body, 1_000_000)

	var request Req
	err := json.NewDecoder(limitedReader).Decode(&request)
	if err != nil {
		resWriter.WriteHeader(http.StatusBadRequest)
		writeErrorText(resWriter, "decoding JSON", err)
		return
	}

	err = request.Validate()
	if err != nil {
		resWriter.WriteHeader(http.StatusBadRequest)
		writeErrorText(resWriter, "validating request", err)
		return
	}

	response, err := w.fn(ctx, request)
	if err != nil {
		resWriter.WriteHeader(http.StatusInternalServerError)
		writeErrorText(resWriter, "running handler", err)
		return
	}

	rawJSON, err := json.Marshal(response)
	if err != nil {
		resWriter.WriteHeader(http.StatusInternalServerError)
		writeErrorText(resWriter, "encoding JSON", err)
		return
	}

	resWriter.Header().Add("Content-Type", "application/json")
	_, _ = resWriter.Write(rawJSON)
}

func writeErrorText(w http.ResponseWriter, text string, err error) {
	buf := bytes.NewBufferString(text)

	buf.WriteString(": ")
	buf.WriteString(err.Error())
	buf.WriteByte('\n')

	w.Write(buf.Bytes())
}

package dbwrapper

import (
	"context"
	"route256/libs/metrics"
	"route256/libs/postgres/transactor"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/opentracing/opentracing-go"
)

type dbWrapper struct {
	db transactor.DB
}

func NewWrapper(db transactor.DB) *dbWrapper {
	return &dbWrapper{
		db: db,
	}
}

func (w *dbWrapper) Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "postgres Query")
	defer span.Finish()

	span.SetTag("query", query)

	timeStart := time.Now()
	rows, err := w.db.Query(ctx, query, args...)
	elapsedTime := time.Since(timeStart)

	metrics.HistogramDBResponseTime.WithLabelValues(query).Observe(elapsedTime.Seconds())

	return rows, err
}

func (w *dbWrapper) QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	span, ctx := opentracing.StartSpanFromContext(ctx, "postgres QueryRow")
	defer span.Finish()

	span.SetTag("query", query)

	timeStart := time.Now()
	row := w.db.QueryRow(ctx, query, args...)
	elapsedTime := time.Since(timeStart)

	metrics.HistogramDBResponseTime.WithLabelValues(query).Observe(elapsedTime.Seconds())

	return row
}

func (w *dbWrapper) Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "postgres Exec")
	defer span.Finish()

	span.SetTag("query", query)

	timeStart := time.Now()
	commandTag, err := w.db.Exec(ctx, query, args...)
	elapsedTime := time.Since(timeStart)

	metrics.HistogramDBResponseTime.WithLabelValues(query).Observe(elapsedTime.Seconds())

	return commandTag, err
}

func (w *dbWrapper) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
	return w.db.BeginTx(ctx, txOptions)
}

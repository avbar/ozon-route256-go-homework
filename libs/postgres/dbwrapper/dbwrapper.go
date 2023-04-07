package dbwrapper

import (
	"context"
	"route256/libs/postgres/transactor"

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

	return w.db.Query(ctx, query, args...)
}

func (w *dbWrapper) QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	span, ctx := opentracing.StartSpanFromContext(ctx, "postgres QueryRow")
	defer span.Finish()

	span.SetTag("query", query)

	return w.db.QueryRow(ctx, query, args...)
}

func (w *dbWrapper) Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "postgres Exec")
	defer span.Finish()

	span.SetTag("query", query)

	return w.db.Exec(ctx, query, args...)
}

func (w *dbWrapper) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
	return w.db.BeginTx(ctx, txOptions)
}

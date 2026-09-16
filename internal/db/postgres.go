package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgxRowsWrapper struct {
	rows pgx.Rows
}

func (w *pgxRowsWrapper) Columns() ([]string, error) {
	fields := w.rows.FieldDescriptions()
	names := make([]string, len(fields))
	for i, f := range fields {
		names[i] = f.Name
	}
	return names, nil
}
func (w *pgxRowsWrapper) Next() bool {
	return w.rows.Next()
}

func (w *pgxRowsWrapper) Scan(dest ...any) error {
	return w.rows.Scan(dest...)
}

func (w *pgxRowsWrapper) Close() error {
	w.rows.Close()
	return nil
}
func (w *pgxRowsWrapper) Err() error {
	return w.rows.Err()
}

type pgxClient struct {
	pool *pgxpool.Pool
}
func (c *pgxClient) Query(ctx context.Context, sql string, args ...any) (Rows, error) {
	rows, err := c.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	return &pgxRowsWrapper{rows: rows}, nil
}
func (c *pgxClient) Ping(ctx context.Context) error {
	return c.pool.Ping(ctx)
}
func (c *pgxClient) Close() error {
	c.pool.Close()
	return nil
}

func NewPgxClient(ctx context.Context, dsn string) (Client, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}
	return &pgxClient{pool: pool}, nil
}

func NewPgxConnector() Connector {
	return func(ctx context.Context, cfg Config) (Client, error) {
		return NewPgxClient(ctx, cfg.ConnectionString())
	}
}
package postgres

import (
	"app-dev/internal/domain"
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

type Client struct {
	pool *pgxpool.Pool
}

func (c *Client) Query(ctx context.Context, sql string, args ...any) (domain.Rows, error) {
	rows, err := c.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	return &pgxRowsWrapper{rows: rows}, nil
}

func (c *Client) Ping(ctx context.Context) error {
	return c.pool.Ping(ctx)
}

func (c *Client) Close() error {
	c.pool.Close()
	return nil
}

func NewClient(ctx context.Context, dsn string) (domain.Client, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}
	return &Client{pool: pool}, nil
}

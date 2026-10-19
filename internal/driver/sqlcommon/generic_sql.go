package sqlcommon

import (
	"app-dev/internal/domain"
	"context"
	"database/sql"
)

// SQLClient wraps *sql.DB to satisfy domain.Client and domain.Querier.
type SQLClient struct {
	DB *sql.DB
}

func (s *SQLClient) Query(ctx context.Context, query string, args ...any) (domain.Rows, error) {
	rows, err := s.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return &SQLRowsWrapper{Rows: rows}, nil
}

func (s *SQLClient) Ping(ctx context.Context) error {
	return s.DB.PingContext(ctx)
}

func (s *SQLClient) Close() error {
	return s.DB.Close()
}

// SQLRowsWrapper wraps *sql.Rows to satisfy domain.Rows.
type SQLRowsWrapper struct {
	Rows *sql.Rows
}

func (r *SQLRowsWrapper) Columns() ([]string, error) {
	return r.Rows.Columns()
}

func (r *SQLRowsWrapper) Next() bool {
	return r.Rows.Next()
}

func (r *SQLRowsWrapper) Scan(dest ...any) error {
	return r.Rows.Scan(dest...)
}

func (r *SQLRowsWrapper) Close() error {
	return r.Rows.Close()
}

func (r *SQLRowsWrapper) Err() error {
	return r.Rows.Err()
}

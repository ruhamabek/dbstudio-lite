package db

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	_ "modernc.org/sqlite"
)

 type SQLClient struct {
	db *sql.DB
}

func (s *SQLClient) Query(ctx context.Context, query string, args ...any) (Rows, error) {
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return &sqlRowsWrapper{rows: rows}, nil
}

func (s *SQLClient) Ping(ctx context.Context) error {
	return s.db.PingContext(ctx)
}

func (s *SQLClient) Close() error {
	return s.db.Close()
}

 type sqlRowsWrapper struct {
	rows *sql.Rows
}

func (r *sqlRowsWrapper) Columns() ([]string, error) {
	return r.rows.Columns()
}

func (r *sqlRowsWrapper) Next() bool {
	return r.rows.Next()
}

func (r *sqlRowsWrapper) Scan(dest ...any) error {
	return r.rows.Scan(dest...)
}

func (r *sqlRowsWrapper) Close() error {
	return r.rows.Close()
}

func (r *sqlRowsWrapper) Err() error {
	return r.rows.Err()
}

 func NewMySQLClient(ctx context.Context, dsn string) (Client, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open mysql connection: %w", err)
	}

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("failed to ping mysql database: %w", err)
	}

	return &SQLClient{db: db}, nil
}

 func NewSQLiteClient(ctx context.Context, path string) (Client, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite database: %w", err)
	}

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("failed to ping sqlite database: %w", err)
	}

	return &SQLClient{db: db}, nil
}
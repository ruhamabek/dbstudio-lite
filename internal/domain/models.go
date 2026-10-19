package domain

import (
	"context"
	"time"
)

// Rows abstracts database row iteration.
type Rows interface {
	Columns() ([]string, error)
	Next() bool
	Scan(dest ...any) error
	Close() error
	Err() error
}

// Querier represents any executor capable of running parameterized queries.
type Querier interface {
	Query(ctx context.Context, sql string, args ...any) (Rows, error)
}

// Client represents a connected database session.
type Client interface {
	Querier
	Ping(ctx context.Context) error
	Close() error
}

// Connector is a factory function for instantiating a Client from Config.
type Connector func(ctx context.Context, cfg Config) (Client, error)

// SchemaInspector abstracts table and column metadata inspection.
type SchemaInspector interface {
	ListTables(ctx context.Context) ([]Table, error)
	GetColumns(ctx context.Context, schema, table string) ([]Column, error)
	GetPrimaryKeys(ctx context.Context, schema, table string) ([]string, error)
}

// ConnectionStore abstracts profile persistence.
type ConnectionStore interface {
	List() ([]SavedConnection, error)
	Save(conn SavedConnection) (SavedConnection, error)
	Delete(id string) error
}

// Table represents a database table in a schema.
type Table struct {
	Schema string
	Name   string
}

// Column represents column metadata.
type Column struct {
	Name       string
	DataType   string
	IsNullable bool
	DefaultVal *string
}

// QueryResult represents the tabular output of a SQL query.
type QueryResult struct {
	Columns      []string `json:"columns"`
	Rows         [][]any  `json:"rows"`
	RowsAffected int64    `json:"rowsAffected"`
	DurationMs   int64    `json:"durationMs"`
}

// SavedConnection represents a persisted connection profile.
type SavedConnection struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"createdAt"`
}

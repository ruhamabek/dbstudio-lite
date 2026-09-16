package db

import (
	"context"
	"fmt"
)

type Rows interface {
	Columns() ([]string, error)
	Next() bool
	Scan(dest ...any) error
	Close() error
	Err() error
}

type Column struct {
	Name       string
	DataType   string
	IsNullable bool
	DefaultVal *string
}

type Querier interface {
	Query(ctx context.Context, sql string, args ...any) (Rows, error)
}

type Table struct {
	Schema string
	Name   string
}

 type SchemaInspector interface {
	ListTables(ctx context.Context) ([]Table, error)
	GetColumns(ctx context.Context, schema, table string) ([]Column, error)
	GetPrimaryKeys(ctx context.Context, schema, table string) ([]string, error)
}

 func NewInspector(driver DriverType, db Querier) SchemaInspector {
	switch driver {
	case DriverMySQL:
		return &MySQLInspector{db: db}
	case DriverSQLite:
		return &SQLiteInspector{db: db}
	default:
		return &PostgresInspector{db: db}
	}
}


type PostgresInspector struct {
	db Querier
}

func (ins *PostgresInspector) ListTables(ctx context.Context) ([]Table, error) {
	query := `
		SELECT table_schema, table_name
		FROM information_schema.tables
		WHERE table_schema NOT IN ('pg_catalog', 'information_schema')
		ORDER BY table_schema, table_name;
	`
	rows, err := ins.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query tables: %w", err)
	}
	defer rows.Close()

	var tables []Table
	for rows.Next() {
		var t Table
		if err := rows.Scan(&t.Schema, &t.Name); err != nil {
			return nil, fmt.Errorf("failed to scan table row: %w", err)
		}
		tables = append(tables, t)
	}
	return tables, rows.Err()
}

func (ins *PostgresInspector) GetColumns(ctx context.Context, schema, table string) ([]Column, error) {
	query := `
		SELECT column_name, data_type, is_nullable, column_default
		FROM information_schema.columns
		WHERE table_schema = $1 AND table_name = $2
		ORDER BY ordinal_position;
	`
	rows, err := ins.db.Query(ctx, query, schema, table)
	if err != nil {
		return nil, fmt.Errorf("failed to query columns: %w", err)
	}
	defer rows.Close()

	var columns []Column
	for rows.Next() {
		var col Column
		var isNullableStr string
		if err := rows.Scan(&col.Name, &col.DataType, &isNullableStr, &col.DefaultVal); err != nil {
			return nil, fmt.Errorf("failed to scan column row: %w", err)
		}
		col.IsNullable = (isNullableStr == "YES")
		columns = append(columns, col)
	}
	return columns, rows.Err()
}

func (ins *PostgresInspector) GetPrimaryKeys(ctx context.Context, schema, table string) ([]string, error) {
	query := `
		SELECT kcu.column_name
		FROM information_schema.table_constraints tc
		JOIN information_schema.key_column_usage kcu
		  ON tc.constraint_name = kcu.constraint_name
		  AND tc.table_schema = kcu.table_schema
		WHERE tc.constraint_type = 'PRIMARY KEY'
		  AND tc.table_schema = $1
		  AND tc.table_name = $2
		ORDER BY kcu.ordinal_position;
	`
	rows, err := ins.db.Query(ctx, query, schema, table)
	if err != nil {
		return nil, fmt.Errorf("failed to get primary keys: %w", err)
	}
	defer rows.Close()

	var pks []string
	for rows.Next() {
		var pk string
		if err := rows.Scan(&pk); err != nil {
			return nil, fmt.Errorf("failed to scan primary key: %w", err)
		}
		pks = append(pks, pk)
	}
	return pks, rows.Err()
}


type MySQLInspector struct {
	db Querier
}

func (ins *MySQLInspector) ListTables(ctx context.Context) ([]Table, error) {
	query := `
		SELECT table_schema, table_name
		FROM information_schema.tables
		WHERE table_schema NOT IN ('information_schema', 'mysql', 'performance_schema', 'sys')
		ORDER BY table_schema, table_name;
	`
	rows, err := ins.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query mysql tables: %w", err)
	}
	defer rows.Close()

	var tables []Table
	for rows.Next() {
		var t Table
		if err := rows.Scan(&t.Schema, &t.Name); err != nil {
			return nil, fmt.Errorf("failed to scan mysql table row: %w", err)
		}
		tables = append(tables, t)
	}
	return tables, rows.Err()
}

func (ins *MySQLInspector) GetColumns(ctx context.Context, schema, table string) ([]Column, error) {
	query := `
		SELECT column_name, data_type, is_nullable, column_default
		FROM information_schema.columns
		WHERE table_schema = ? AND table_name = ?
		ORDER BY ordinal_position;
	`
	rows, err := ins.db.Query(ctx, query, schema, table)
	if err != nil {
		return nil, fmt.Errorf("failed to query mysql columns: %w", err)
	}
	defer rows.Close()

	var columns []Column
	for rows.Next() {
		var col Column
		var isNullableStr string
		if err := rows.Scan(&col.Name, &col.DataType, &isNullableStr, &col.DefaultVal); err != nil {
			return nil, fmt.Errorf("failed to scan mysql column row: %w", err)
		}
		col.IsNullable = (isNullableStr == "YES")
		columns = append(columns, col)
	}
	return columns, rows.Err()
}

func (ins *MySQLInspector) GetPrimaryKeys(ctx context.Context, schema, table string) ([]string, error) {
	query := `
		SELECT column_name
		FROM information_schema.key_column_usage
		WHERE table_schema = ? AND table_name = ? AND constraint_name = 'PRIMARY'
		ORDER BY ordinal_position;
	`
	rows, err := ins.db.Query(ctx, query, schema, table)
	if err != nil {
		return nil, fmt.Errorf("failed to get mysql primary keys: %w", err)
	}
	defer rows.Close()

	var pks []string
	for rows.Next() {
		var pk string
		if err := rows.Scan(&pk); err != nil {
			return nil, fmt.Errorf("failed to scan mysql primary key: %w", err)
		}
		pks = append(pks, pk)
	}
	return pks, rows.Err()
}


type SQLiteInspector struct {
	db Querier
}

func (ins *SQLiteInspector) ListTables(ctx context.Context) ([]Table, error) {
	query := `
		SELECT 'main' AS table_schema, name AS table_name
		FROM sqlite_master
		WHERE type = 'table' AND name NOT LIKE 'sqlite_%'
		ORDER BY name;
	`
	rows, err := ins.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query sqlite tables: %w", err)
	}
	defer rows.Close()

	var tables []Table
	for rows.Next() {
		var t Table
		if err := rows.Scan(&t.Schema, &t.Name); err != nil {
			return nil, fmt.Errorf("failed to scan sqlite table row: %w", err)
		}
		tables = append(tables, t)
	}
	return tables, rows.Err()
}

func (ins *SQLiteInspector) GetColumns(ctx context.Context, schema, table string) ([]Column, error) {
	query := fmt.Sprintf("PRAGMA table_info(%q);", table)
	rows, err := ins.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query sqlite columns: %w", err)
	}
	defer rows.Close()

	var columns []Column
	for rows.Next() {
		var cid int64
		var name, dataType string
		var notNull, pk int64
		var dfltVal *string

		if err := rows.Scan(&cid, &name, &dataType, &notNull, &dfltVal, &pk); err != nil {
			return nil, fmt.Errorf("failed to scan sqlite column row: %w", err)
		}

		columns = append(columns, Column{
			Name:       name,
			DataType:   dataType,
			IsNullable: notNull == 0,
			DefaultVal: dfltVal,
		})
	}
	return columns, rows.Err()
}

func (ins *SQLiteInspector) GetPrimaryKeys(ctx context.Context, schema, table string) ([]string, error) {
	query := fmt.Sprintf("PRAGMA table_info(%q);", table)
	rows, err := ins.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query sqlite primary keys: %w", err)
	}
	defer rows.Close()

	var pks []string
	for rows.Next() {
		var cid int64
		var name, dataType string
		var notNull, pk int64
		var dfltVal *string

		if err := rows.Scan(&cid, &name, &dataType, &notNull, &dfltVal, &pk); err != nil {
			return nil, fmt.Errorf("failed to scan sqlite pk row: %w", err)
		}

		if pk > 0 {
			pks = append(pks, name)
		}
	}
	return pks, rows.Err()
}
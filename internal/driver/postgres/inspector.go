package postgres

import (
	"app-dev/internal/domain"
	"context"
	"fmt"
)

type Inspector struct {
	db domain.Querier
}

func NewInspector(db domain.Querier) *Inspector {
	return &Inspector{db: db}
}

func (ins *Inspector) ListTables(ctx context.Context) ([]domain.Table, error) {
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

	var tables []domain.Table
	for rows.Next() {
		var t domain.Table
		if err := rows.Scan(&t.Schema, &t.Name); err != nil {
			return nil, fmt.Errorf("failed to scan table row: %w", err)
		}
		tables = append(tables, t)
	}
	return tables, rows.Err()
}

func (ins *Inspector) GetColumns(ctx context.Context, schema, table string) ([]domain.Column, error) {
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

	var columns []domain.Column
	for rows.Next() {
		var col domain.Column
		var isNullableStr string
		if err := rows.Scan(&col.Name, &col.DataType, &isNullableStr, &col.DefaultVal); err != nil {
			return nil, fmt.Errorf("failed to scan column row: %w", err)
		}
		col.IsNullable = (isNullableStr == "YES")
		columns = append(columns, col)
	}
	return columns, rows.Err()
}

func (ins *Inspector) GetPrimaryKeys(ctx context.Context, schema, table string) ([]string, error) {
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

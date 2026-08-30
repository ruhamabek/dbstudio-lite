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

type Inspector struct {
	db Querier
}

func NewInspector(db Querier) *Inspector{
	return &Inspector{db: db}
}


func (ins *Inspector) ListTables(ctx context.Context) ([]Table, error) {
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
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating table rows: %w", err)
	}
	return tables, nil
}

func (ins *Inspector) GetColumns(ctx context.Context, schema, table string)([]Column, error){
	query := `
		SELECT column_name, data_type, is_nullable, column_default
		FROM information_schema.columns
		WHERE table_schema = $1 AND table_name = $2
		ORDER BY ordinal_position;
	`

    rows, err := ins.db.Query(ctx, query, schema, table)
	if err != nil {
		return nil, fmt.Errorf("Failed to query columns: %v", err)
	}

	defer rows.Close()

	var columns []Column
	for rows.Next(){
		var col Column
		var isNullableStr string
		if err := rows.Scan(&col.Name, &col.DataType, &isNullableStr, &col.DefaultVal); err != nil {
			return nil, fmt.Errorf("Failed to scan column row: %v", err)
		}

		col.IsNullable = (isNullableStr =="YES")

		columns = append(columns, col)
	}

	if err := rows.Err(); err != nil {
		return nil,fmt.Errorf("error iterating the columns: %v", err)
	}

	return columns,nil
}

func (ins *Inspector) GetPrimaryKeys(ctx context.Context, schema, table string)([]string, error){
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
		return nil, fmt.Errorf("Failed to get primary keys: %v", err)
	}

	defer rows.Close()

	pks := []string{}

	for rows.Next(){
		var pk string
		if err:= rows.Scan(&pk); err !=nil {
			return nil, fmt.Errorf("failed to get primary keys: %v", err)
		}

		pks = append(pks, pk)
	}

	return pks, nil
}
package sqlite

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

	var tables []domain.Table
	for rows.Next() {
		var t domain.Table
		if err := rows.Scan(&t.Schema, &t.Name); err != nil {
			return nil, fmt.Errorf("failed to scan sqlite table row: %w", err)
		}
		tables = append(tables, t)
	}
	return tables, rows.Err()
}

func (ins *Inspector) GetColumns(ctx context.Context, schema, table string) ([]domain.Column, error) {
	query := fmt.Sprintf("PRAGMA table_info(%q);", table)
	rows, err := ins.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query sqlite columns: %w", err)
	}
	defer rows.Close()

	var columns []domain.Column
	for rows.Next() {
		var cid int64
		var name, dataType string
		var notNull, pk int64
		var dfltVal *string

		if err := rows.Scan(&cid, &name, &dataType, &notNull, &dfltVal, &pk); err != nil {
			return nil, fmt.Errorf("failed to scan sqlite column row: %w", err)
		}

		columns = append(columns, domain.Column{
			Name:       name,
			DataType:   dataType,
			IsNullable: notNull == 0,
			DefaultVal: dfltVal,
		})
	}
	return columns, rows.Err()
}

func (ins *Inspector) GetPrimaryKeys(ctx context.Context, schema, table string) ([]string, error) {
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

package service

import (
	"app-dev/internal/domain"
	"context"
	"fmt"
	"time"
)

type Runner struct {
	db domain.Querier
}

func NewRunner(db domain.Querier) *Runner {
	return &Runner{db: db}
}

func (r *Runner) Execute(ctx context.Context, sql string, args ...any) (*domain.QueryResult, error) {
	start := time.Now()

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("query execution failed: %w", err)
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve columns: %w", err)
	}

	result := &domain.QueryResult{
		Columns: cols,
		Rows:    [][]any{},
	}

	for rows.Next() {
		rowValues := make([]any, len(cols))
		scanArgs := make([]any, len(cols))

		for i := range rowValues {
			scanArgs[i] = &rowValues[i]
		}

		if err := rows.Scan(scanArgs...); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		result.Rows = append(result.Rows, rowValues)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	result.DurationMs = int64(time.Since(start).Milliseconds())

	return result, nil
}

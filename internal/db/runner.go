package db

import (
	"context"
	"fmt"
	"time"
)

type QueryResult struct {
	Columns      []string `json:"columns"`
	Rows         [][]any  `json:"rows"`
	RowsAffected int64    `json:"rowsAffected"`
	DurationMs   int64    `json:"durationMs"`
}

type Runner struct {
	db Querier
}

func NewRunner(db Querier) *Runner {
	return &Runner{db: db}
}

func (r *Runner) Execute(ctx context.Context, sql string, args ...any)(* QueryResult, error){
	start := time.Now()

	rows, err := r.db.Query(ctx, sql, args...)

	if err != nil {
		 return nil, fmt.Errorf("Query execution failed: %v", err)
	}

	defer rows.Close()

	cols, err := rows.Columns()

	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve columns: %v", err)
	}

	result := &QueryResult{
		Columns: cols,
		Rows: [][]any{},
	}

	for rows.Next(){
		rowValues := make([]any, len(cols))
		scanArgs := make([]any, len(cols))

		for i := range rowValues {
			 scanArgs[i] = &rowValues[i]
		}

		if err := rows.Scan(scanArgs...); err != nil {
			 return nil, fmt.Errorf("failed to scan row: %v", err)
		}

		result.Rows = append(result.Rows, rowValues)
	}

	if err := rows.Err(); err != nil {
		 return nil, fmt.Errorf("errror during rows iteration: %v", err)
	}

	result.DurationMs = int64(time.Since(start).Milliseconds())

	return result, nil
}
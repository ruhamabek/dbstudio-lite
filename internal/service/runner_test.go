package service_test

import (
	"app-dev/internal/domain"
	"app-dev/internal/service"
	"context"
	"fmt"
	"testing"
)

type FakeRunnerRows struct {
	cols []string
	data [][]any
	idx  int
}

func (r *FakeRunnerRows) Columns() ([]string, error) {
	return r.cols, nil
}

func (r *FakeRunnerRows) Next() bool {
	r.idx++
	return r.idx <= len(r.data)
}

func (r *FakeRunnerRows) Scan(dest ...any) error {
	if r.idx > len(r.data) {
		return fmt.Errorf("out of range")
	}
	row := r.data[r.idx-1]
	for i, d := range dest {
		if p, ok := d.(*any); ok {
			*p = row[i]
		}
	}
	return nil
}

func (r *FakeRunnerRows) Close() error { return nil }
func (r *FakeRunnerRows) Err() error   { return nil }

type FakeRunnerQuerier struct {
	rowsToReturn *FakeRunnerRows
}

func (q *FakeRunnerQuerier) Query(ctx context.Context, sql string, args ...any) (domain.Rows, error) {
	return q.rowsToReturn, nil
}

func TestRunner_Execute(t *testing.T) {
	t.Run("executes query and returns dynamic tabular results", func(t *testing.T) {
		fakeRows := &FakeRunnerRows{
			cols: []string{"id", "title", "completed"},
			data: [][]any{
				{1, "Setup project", true},
				{2, "Implement TDD", false},
			},
		}

		runner := service.NewRunner(&FakeRunnerQuerier{rowsToReturn: fakeRows})
		result, err := runner.Execute(context.Background(), "SELECT * FROM tasks")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(result.Columns) != 3 {
			t.Errorf("got %d columns, want 3", len(result.Columns))
		}

		if len(result.Rows) != 2 {
			t.Errorf("got %d rows, want 2", len(result.Rows))
		}

		if result.Rows[0][1] != "Setup project" {
			t.Errorf("unexpected cell value: got %v, want 'Setup project'", result.Rows[0][1])
		}
	})
}

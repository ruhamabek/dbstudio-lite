package db_test

import (
	"context"
	"fmt"
	"testing"
	"app-dev/internal/db"
)

type DynamicFakeRows struct {
	cols []string
	data [][]any
	idx  int
}

func (d *DynamicFakeRows) Columns() ([]string, error){
	 return d.cols,nil
}

func (d *DynamicFakeRows) Next() bool{
	d.idx++
	return d.idx <= len(d.data)
}

func (d *DynamicFakeRows) Scan(dest ...any) error{
       if d.idx > len(d.data){
		return fmt.Errorf("out of range")
	   }

	   row := d.data[d.idx - 1]
	   for i, val := range row {
		  if ptr, ok := dest[i].(*any); ok {
			  *ptr = val
		  }
	   }
	
	   return nil
}


func (d *DynamicFakeRows) Close() error { return nil }
func (d *DynamicFakeRows) Err() error   { return nil }

type DynamicFakeQuerier struct{
	rowsToReturn *DynamicFakeRows
}

func (dq *DynamicFakeQuerier) Query(ctx context.Context, sql string, args ...any)(db.Rows, error){
	return dq.rowsToReturn, nil
}

func TestRunner_Execute(t *testing.T) {
	t.Run("executes query and returns dynamic tabular results", func(t *testing.T) {
		  fakeRows := &DynamicFakeRows{
			cols : []string{"id", "username", "is_admin"},
			data: [][]any{
				{1, "alice", true},
				{2, "bob", false},
			},
		  }
		  fakeDB := &DynamicFakeQuerier{
			 rowsToReturn: fakeRows,
		  }

		  runner := db.NewRunner(fakeDB)
		  
		  result, err := runner.Execute(context.Background(), "SELECT id, username, is_admin FROM users")
         
		  if err != nil {
			   t.Fatalf("unexpected error: %v", err)
		  }

		  wantCols := []string{"id", "username", "is_admin"}

		  if len(result.Columns) != len(wantCols) {
			  t.Fatalf("got %v, want %v", len(result.Columns), len(wantCols))
		  }

		  for i := range wantCols{
			 if result.Columns[i] != wantCols[i] {
				t.Errorf("col %v, got %v, want %v", i, result.Columns[i], wantCols[i])
			 }
		  }

		 if len(result.Rows) != 2 {
			t.Fatalf("got %d rows, want 2", len(result.Rows))
		}
		 if result.Rows[0][1] != "alice" || result.Rows[1][1] != "bob" {
			t.Errorf("unexpected rows data: %+v", result.Rows)
		 }

		 if result.DurationMs < 0 {
			t.Errorf("expected duration >= 0, got %d", result.DurationMs)
		 }
	})
}
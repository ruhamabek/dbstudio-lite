package db_test

import (
	"context"
	"fmt"
	"testing"
    "app-dev/internal/db"

)


type FakeRows struct {
	data [][]any
	idx int
}

func (f *FakeRows) Next() bool{
    f.idx++
	return f.idx <= len(f.data)
}

func (f *FakeRows) Scan(dest ...any) error {
       if f.idx > len(f.data){
		  return fmt.Errorf("Out of range")
	   }

	   row := f.data[f.idx-1]
	   for i, d := range dest {
		  switch target := d.(type){
		   case *string:
			if row[i] != nil {
			*target = row[i].(string)
		  }

		 case **string:
			if row[i] == nil {
				*target = nil
			} else {
				*target = row[i].(*string)
			}
		
	   }
	}
	   return nil
}

func (f *FakeRows) Close()error { return nil }
func (f *FakeRows) Err() error   { return nil }
func (f *FakeRows) Columns() ([]string, error) { return nil, nil }

type FakeQuerier struct {
	rowsToReturn *FakeRows
	recordedSQL  string
}



func (fq *FakeQuerier) Query(ctx context.Context, query string, args ...any)(db.Rows, error){
	  fq.recordedSQL = query
	  return fq.rowsToReturn, nil
}

func TestInspector_ListTables(t *testing.T) {
	t.Run("return list of tables from databases", func(t *testing.T) {
		  fakeRows := &FakeRows{
			data: [][]any{
				{"public", "users"},
				{"public", "posts"},
				{"analytics", "events"},
			},
	}
	fakeDB := &FakeQuerier{rowsToReturn: fakeRows}

	inspector := db.NewInspector(fakeDB)

	tables, err := inspector.ListTables(context.Background())

	if err != nil {
		t.Fatalf("unexpected error :%v", err)
	}

	want := []db.Table{
		{Schema: "public", Name: "users"},
		{Schema: "public", Name: "posts"},
		{Schema: "analytics", Name: "events"},
	}
	
	if len(tables) != len(want){
			t.Fatalf("got %d tables, want %d", len(tables), len(want))
	}

	for i := range want {
			if tables[i] != want[i] {
				t.Errorf("at index %d: got %+v, want %+v", i, tables[i], want[i])
			}
		}
	})
}

func TestInspector_GetColumns(t *testing.T) {
	t.Run("Returns table columns for a specific table", func(t *testing.T) {
		  defaultVal := "now()"
		  fakeRows := &FakeRows{
			data: [][]any{
				{"id", "integer", "NO", nil},
				{"email", "varchar", "NO", nil},
				{"created_at", "timestamptz", "YES", &defaultVal},
			},
		}
		fakeDB := &FakeQuerier{
			rowsToReturn: fakeRows,
		}

		inspector := db.NewInspector(fakeDB)

		cols,err := inspector.GetColumns(context.Background(), "public", "users")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(cols) != 3 {
			t.Fatalf("got %v, want 3", len(cols))
		}

		if cols[0].Name != "id" || cols[0].DataType != "integer" || cols[0].IsNullable != false {
			t.Errorf("unexpected column 0: %+v", cols[0])
		}
 
		if cols[2].Name != "created_at" || !cols[2].IsNullable || cols[2].DefaultVal == nil || *cols[2].DefaultVal != "now()" {
			t.Errorf("unexpected column 2: %+v", cols[2])
		}


	})
}

func TestInspector_GetPrimaryKeys(t *testing.T) {
	t.Run("returns primary key", func(t *testing.T) {
		 fakeRows := &FakeRows{
			data: [][]any{
				{"order_id"},
				{"item_id"},
			},
		 }

		 fakeDB := &FakeQuerier{rowsToReturn: fakeRows}
		 inspector := db.NewInspector(fakeDB)
		 pks, err := inspector.GetPrimaryKeys(context.Background(), "public", "order_items")
		 if err != nil {
			  t.Fatalf("unexpected error: %v", err)
		 }

		 want := []string{"order_id", "item_id"}

		 if len(pks) != len(want){
			 t.Fatalf("got %v, want %v", pks, want)
		 }

		 for i := range want {
			  if pks[i] != want[i]{
				t.Errorf("at index %d: got %q, want %q", i, pks[i], want[i])
			  }
		 }

	})

	t.Run("return empty slice when table has no primary key", func(t *testing.T) {
		   fakeRows := &FakeRows{
			   data: [][]any{},
		   }

		   fakeDB := &FakeQuerier{rowsToReturn : fakeRows}

		   inspector := db.NewInspector(fakeDB)

		   pks, err := inspector.GetPrimaryKeys(context.Background(), "public", "logs")

		   if err != nil {
			t.Fatalf("failed to get primary key: %v", err)
		   }

		   if len(pks) != 0 {
			    t.Errorf("got %v, want 0", len(pks))
		   }
	})
}
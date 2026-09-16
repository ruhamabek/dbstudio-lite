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
	if f.idx > len(f.data) {
		return fmt.Errorf("out of range")
	}

	row := f.data[f.idx-1]
	for i, d := range dest {
		if i >= len(row) || row[i] == nil {
			continue
		}
		switch target := d.(type) {
		case *string:
			*target = fmt.Sprintf("%v", row[i])
		case **string:
			if strPtr, ok := row[i].(*string); ok {
				*target = strPtr
			} else {
				val := fmt.Sprintf("%v", row[i])
				*target = &val
			}
		case *int:
			if v, ok := row[i].(int); ok {
				*target = v
			} else if v, ok := row[i].(int64); ok {
				*target = int(v)
			}
		case *int64:
			if v, ok := row[i].(int64); ok {
				*target = v
			} else if v, ok := row[i].(int); ok {
				*target = int64(v)
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


func (fq *FakeQuerier) Query(ctx context.Context, query string, args ...any) (db.Rows, error) {
	fq.recordedSQL = query
	if fq.rowsToReturn != nil {
		fq.rowsToReturn.idx = 0  
	}
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

	inspector := db.NewInspector(db.DriverPostgres, fakeDB)

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

		inspector := db.NewInspector(db.DriverPostgres, fakeDB)

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
		 inspector := db.NewInspector(db.DriverPostgres, fakeDB)
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

		   inspector := db.NewInspector(db.DriverSQLite, fakeDB)

		   pks, err := inspector.GetPrimaryKeys(context.Background(), "public", "logs")

		   if err != nil {
			t.Fatalf("failed to get primary key: %v", err)
		   }

		   if len(pks) != 0 {
			    t.Errorf("got %v, want 0", len(pks))
		   }
	})
}

func TestMySQLInspector(t *testing.T) {
	t.Run("lists tables for current mysql database", func(t *testing.T) {
		fakeDB := &FakeQuerier{
			rowsToReturn: &FakeRows{
				data: [][]any{
					{"mydb", "users"},
					{"mydb", "orders"},
				},
			},
		}

		inspector := db.NewInspector(db.DriverMySQL, fakeDB)
		tables, err := inspector.ListTables(context.Background())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(tables) != 2 || tables[0].Name != "users" || tables[1].Name != "orders" {
			t.Errorf("unexpected tables: %+v", tables)
		}
	})

	t.Run("gets primary keys in mysql using ? placeholders", func(t *testing.T) {
		fakeDB := &FakeQuerier{
			rowsToReturn: &FakeRows{
				data: [][]any{
					{"id"},
				},
			},
		}

		inspector := db.NewInspector(db.DriverMySQL, fakeDB)
		pks, err := inspector.GetPrimaryKeys(context.Background(), "mydb", "users")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(pks) != 1 || pks[0] != "id" {
			t.Errorf("unexpected pks: %+v", pks)
		}
	})
}

func TestSQLiteInspector(t *testing.T) {
	t.Run("lists tables from sqlite_master", func(t *testing.T) {
		fakeDB := &FakeQuerier{
			rowsToReturn: &FakeRows{
				data: [][]any{
					{"main", "products"},
					{"main", "categories"},
				},
			},
		}

		inspector := db.NewInspector(db.DriverSQLite, fakeDB)
		tables, err := inspector.ListTables(context.Background())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(tables) != 2 || tables[0].Name != "products" {
			t.Errorf("unexpected tables: %+v", tables)
		}
	})

	t.Run("gets columns and primary keys using PRAGMA table_info", func(t *testing.T) {
		defaultVal := "active"
		fakeDB := &FakeQuerier{
			rowsToReturn: &FakeRows{
				data: [][]any{
 					{int64(0), "id", "INTEGER", int64(1), nil, int64(1)},          // notnull=1 (IsNullable=false), pk=1
					{int64(1), "status", "TEXT", int64(0), &defaultVal, int64(0)}, // notnull=0 (IsNullable=true), pk=0
				},
			},
		}

		inspector := db.NewInspector(db.DriverSQLite, fakeDB)

 		cols, err := inspector.GetColumns(context.Background(), "main", "products")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(cols) != 2 {
			t.Fatalf("expected 2 columns, got %d", len(cols))
		}
		if cols[0].Name != "id" || cols[0].IsNullable != false {
			t.Errorf("column 0 mismatch: %+v", cols[0])
		}
		if cols[1].Name != "status" || cols[1].IsNullable != true {
			t.Errorf("column 1 mismatch: %+v", cols[1])
		}

 		pks, err := inspector.GetPrimaryKeys(context.Background(), "main", "products")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(pks) != 1 || pks[0] != "id" {
			t.Errorf("expected pk ['id'], got %+v", pks)
		}
	})
}
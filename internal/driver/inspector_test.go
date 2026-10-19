package driver_test

import (
	"app-dev/internal/domain"
	"app-dev/internal/driver"
	"context"
	"fmt"
	"testing"
)

type FakeRows struct {
	data [][]any
	idx  int
}

func (f *FakeRows) Next() bool {
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

func (f *FakeRows) Close() error               { return nil }
func (f *FakeRows) Err() error                 { return nil }
func (f *FakeRows) Columns() ([]string, error) { return nil, nil }

type FakeQuerier struct {
	rowsToReturn *FakeRows
	recordedSQL  string
}

func (fq *FakeQuerier) Query(ctx context.Context, query string, args ...any) (domain.Rows, error) {
	fq.recordedSQL = query
	if fq.rowsToReturn != nil {
		fq.rowsToReturn.idx = 0
	}
	return fq.rowsToReturn, nil
}

func TestPostgresInspector(t *testing.T) {
	t.Run("lists tables from information_schema", func(t *testing.T) {
		fakeRows := &FakeRows{
			data: [][]any{
				{"public", "users"},
				{"public", "orders"},
			},
		}
		fakeDB := &FakeQuerier{rowsToReturn: fakeRows}
		inspector := driver.NewInspector(domain.DriverPostgres, fakeDB)

		tables, err := inspector.ListTables(context.Background())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(tables) != 2 || tables[0].Name != "users" || tables[1].Name != "orders" {
			t.Errorf("unexpected tables: %+v", tables)
		}
	})

	t.Run("gets columns for table", func(t *testing.T) {
		defaultVal := "now()"
		fakeRows := &FakeRows{
			data: [][]any{
				{"id", "integer", "NO", nil},
				{"created_at", "timestamptz", "YES", &defaultVal},
			},
		}
		fakeDB := &FakeQuerier{rowsToReturn: fakeRows}
		inspector := driver.NewInspector(domain.DriverPostgres, fakeDB)

		cols, err := inspector.GetColumns(context.Background(), "public", "users")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(cols) != 2 || cols[0].IsNullable || !cols[1].IsNullable || *cols[1].DefaultVal != "now()" {
			t.Errorf("unexpected columns: %+v", cols)
		}
	})

	t.Run("gets primary keys", func(t *testing.T) {
		fakeRows := &FakeRows{
			data: [][]any{
				{"id"},
			},
		}
		fakeDB := &FakeQuerier{rowsToReturn: fakeRows}
		inspector := driver.NewInspector(domain.DriverPostgres, fakeDB)

		pks, err := inspector.GetPrimaryKeys(context.Background(), "public", "users")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(pks) != 1 || pks[0] != "id" {
			t.Errorf("unexpected primary keys: %+v", pks)
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

		inspector := driver.NewInspector(domain.DriverMySQL, fakeDB)
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

		inspector := driver.NewInspector(domain.DriverMySQL, fakeDB)
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

		inspector := driver.NewInspector(domain.DriverSQLite, fakeDB)
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
					{int64(0), "id", "INTEGER", int64(1), nil, int64(1)},
					{int64(1), "status", "TEXT", int64(0), &defaultVal, int64(0)},
				},
			},
		}

		inspector := driver.NewInspector(domain.DriverSQLite, fakeDB)

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

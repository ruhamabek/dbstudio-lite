package relay_test

import (
	"app-dev/internal/domain"
	"app-dev/internal/relay"
	"context"
	"encoding/json"
	"testing"
)

type FakeInspector struct {
	tables []domain.Table
}



func (f *FakeInspector) ListTables(ctx context.Context) ([]domain.Table, error) {
	return f.tables, nil
}


func (f *FakeInspector) GetColumns(ctx context.Context, schema, table string) ([]domain.Column, error) {
	return []domain.Column{{Name: "id", DataType: "int"}}, nil
}


func (f *FakeInspector) GetPrimaryKeys(ctx context.Context, schema, table string) ([]string, error) {
	return []string{"id"}, nil
}


type FakeQuerier struct{}

type FakeRows struct {
	closed bool
}

func (r *FakeRows) Columns() ([]string, error)  { return []string{"val"}, nil }
func (r *FakeRows) Next() bool                  { return false }
func (r *FakeRows) Scan(dest ...any) error      { return nil }
func (r *FakeRows) Close() error                { r.closed = true; return nil }
func (r *FakeRows) Err() error                  { return nil }

func (q *FakeQuerier) Query(ctx context.Context, sql string, args...any)(domain.Rows, error){
	return &FakeRows{}, nil
}

func TestGateway_Handle(t *testing.T) {
	ctx := context.Background()
	inspector := &FakeInspector{
		tables: []domain.Table{{Schema: "public", Name: "users"}},
	}
	querier := &FakeQuerier{}
	gateway := relay.NewHostGateway(inspector, querier)

	t.Run("Handle returns error for empty SQL", func(t *testing.T) {
		req := relay.Request{
			ID: "req-1",
			Method: relay.MethodListTables,
		}

		resp := gateway.Handle(ctx, relay.PermissionView, req)
		if resp.Error != "" {
			t.Fatalf("expected no error for listing tables, got: %v", resp.Error)
		}

		var tables []domain.Table
			if err := json.Unmarshal(resp.Result, &tables); err != nil || len(tables) != 1 || tables[0].Name != "users" {
				t.Fatalf("expected one table 'users', got: %v", string(resp.Result))
		}
	})

	t.Run("rejects execute_query under view permission", func(t *testing.T) {
		req := relay.Request{
			ID:     "req-2",
			Method: relay.MethodExecuteQuery,
			SQL:    "SELECT * FROM users",
		}
		
		resp := gateway.Handle(ctx, relay.PermissionView, req)
		if resp.Error == "" {
			t.Fatal("expected error under View permission, got success")
		}
	})

	t.Run("rejects mutating query under Read permission", func(t *testing.T) {
		  req := relay.Request{
			ID:     "req-3",
			Method: relay.MethodExecuteQuery,
			SQL:    "DELETE FROM users",
		  }

		  resp := gateway.Handle(ctx, relay.PermissionRead, req)
		  if resp.Error == "" {
			   t.Fatal("expected error for DELETE under Read permission, got success")
		  }
	})

	t.Run("execute valid SELECT query under read permission", func(t *testing.T) {
		 req := relay.Request{
			ID:     "req-4",
			Method: relay.MethodExecuteQuery,
			SQL:    "SELECT * FROM users",
		}

		resp := gateway.Handle(ctx, relay.PermissionRead, req)
		if resp.Error != ""{
			t.Fatalf("expected no error for SELECT under Read permission, got: %v", resp.Error)
		}

		var result domain.QueryResult
		if err := json.Unmarshal(resp.Result, &result); err != nil {
			t.Fatalf("expected valid QueryResult, got: %v", err)
		}
	})

	t.Run("returns error on unknown method", func(t *testing.T) {
		req := relay.Request{
			ID:     "req-5",
			Method: "unknown_action",
		}

		resp := gateway.Handle(ctx, relay.PermissionUpdate, req)
		if resp.Error == "" {
			t.Fatal("expected error for unknown method, got nil")
		}
	})
}

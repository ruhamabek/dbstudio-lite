package relay_test

import (
	"app-dev/internal/relay"
	"context"
	"encoding/json"
	"errors"
	"testing"
)

type FakeCaller struct{
	handler func(req relay.Request) (relay.Response, error)
}

func (f *FakeCaller) Call(ctx context.Context, req relay.Request) (relay.Response, error){
	return f.handler(req)
}


func TestRelayClient(t *testing.T) {
	ctx := context.Background()
    
	t.Run("ListTables fetches tables and unmarshals them", func(t *testing.T) {
		caller := &FakeCaller{
			handler: func(req relay.Request) (relay.Response, error) {
				if req.Method != relay.MethodListTables {
					t.Fatalf("expected method %s, got %s", relay.MethodListTables, req.Method)
				}

				return relay.Response{
					ID:     req.ID,
					Result: json.RawMessage(`[{"Schema":"public","Name":"users"}]`),
				}, nil
			},
		}

		client := relay.NewRelayClient(caller)
		tables, err := client.ListTables(ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(tables) != 1 || tables[0].Name != "users" {
			t.Fatalf("unexpected tables: %+v", tables)
		}
	})

	t.Run("Execute returns QueryResult on success", func(t *testing.T) {
		caller := &FakeCaller{
			handler: func(req relay.Request) (relay.Response, error) {
				if req.Method != relay.MethodExecuteQuery || req.SQL != "SELECT 1" {
					t.Fatalf("unexpected request: %+v", req)
				}
				return relay.Response{
					ID: req.ID,
					Result: json.RawMessage(`{"columns":["?column?"],"rows":[["1"]],"rowsAffected":1,"durationMs":5}`),
				}, nil	
			 
			},
		}

		client := relay.NewRelayClient(caller)
		result, err := client.Execute(ctx, "SELECT 1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result.Columns) != 1 || result.Columns[0] != "?column?" {
			t.Fatalf("unexpected result: %+v", result)
		}
	})
	
	t.Run("returns error when transport fails", func(t *testing.T) {
		caller := &FakeCaller{
			handler: func(req relay.Request) (relay.Response, error) {
				 return relay.Response{}, errors.New("network stream closed")
			},
		}

		client := relay.NewRelayClient(caller)
		_,err := client.ListTables(ctx)
		if err == nil {
			t.Fatal("expected transport error, got nil")
		}
	})

}

package relay_test

import (
	"context"
	"net"
	"testing"
	"time"
	"app-dev/internal/domain"
	"app-dev/internal/relay"
)

func TestStreamRoundTrip(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	serverConn, clientConn := net.Pipe()
	defer serverConn.Close()
	defer clientConn.Close()

	inspector := &FakeInspector{
		tables: []domain.Table{{Schema: "main", Name: "products"}},
	}

	querier := &FakeQuerier{}
	gateway := relay.NewHostGateway(inspector, querier)

	go relay.ServeStream(ctx, serverConn, gateway, relay.PermissionRead)

	caller := relay.NewStreamCaller(clientConn)
	client := relay.NewRelayClient(caller)

	t.Run("successfully lists tables over stream", func(t *testing.T) {
		  tables, err := client.ListTables(ctx)
		  if err != nil {
			 t.Fatalf("unecpeted error over stream: %v", err)
		  }

		  if len(tables) != 1 || tables[0].Name != "products" {
			t.Fatalf("unexpected tables receieved: %+v", tables)
		  }
	})

    t.Run("successfully executes allowed SELECT query over stream", func(t *testing.T) {
		res, err := client.Execute(ctx, "SELECT * FROM products")
		if err != nil {
			t.Fatalf("unexpected query error: %v", err)
		}
		if res == nil {
			t.Fatal("expected non-nil query result")
		}
	})

	t.Run("enforces permission guard across stream", func(t *testing.T) {
		_, err := client.Execute(ctx, "DELETE FROM products")
		if err == nil {
			t.Fatal("expected DELETE to be rejected under Read permission, got nil")
		}
	})
}
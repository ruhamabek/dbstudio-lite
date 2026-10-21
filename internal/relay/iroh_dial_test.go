package relay_test

import (
	"context"
	"net"
	"testing"
	"time"

	"app-dev/internal/domain"
	"app-dev/internal/relay"
)

func TestIrohCrossNetworkDial(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping network test in short mode")
	}

	c, err := net.DialTimeout("tcp", "1.1.1.1:53", 500*time.Millisecond)
	if err != nil {
		t.Skip("skipping live public relay test: no internet connectivity")
	}
	c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()

	inspector := &FakeInspector{
		tables: []domain.Table{{Schema: "public", Name: "remote_products"}},
	}
	querier := &FakeQuerier{}

 	host, ticketStr, err := relay.StartHostSession(ctx, "0.0.0.0:0", inspector, querier, relay.PermissionRead, "store_db")
	if err != nil {
		t.Fatalf("failed to start host session: %v", err)
	}
	defer host.Close()

	parsed, err := relay.ParseTicket(ticketStr)
	if err != nil {
		t.Fatalf("failed to parse ticket: %v", err)
	}

	if parsed.NodeID == "" {
		t.Fatal("expected ticket to contain Iroh NodeID for cross-network reachability")
	}
	t.Logf("Host NodeID: %s, Relays: %v", parsed.NodeID, parsed.RelayURLs)

	// SIMULATE SCENARIO B (DIFFERENT NETWORKS):
	parsed.Addrs = nil
	remoteTicketStr, err := parsed.Encode()
	if err != nil {
		t.Fatalf("failed to encode remote ticket: %v", err)
	}

	client, err := relay.DialTicket(ctx, remoteTicketStr)
	if err != nil {
		t.Fatalf("failed to dial ticket across relay network: %v", err)
	}
	defer client.Close()

	tables, err := client.ListTables(ctx)
	if err != nil {
		t.Fatalf("unexpected error listing tables over relay: %v", err)
	}
	if len(tables) != 1 || tables[0].Name != "remote_products" {
		t.Fatalf("unexpected tables: %+v", tables)
	}

	res, err := client.Execute(ctx, "SELECT * FROM remote_products")
	if err != nil {
		t.Fatalf("unexpected query error over relay: %v", err)
	}
	if res == nil {
		t.Fatal("expected non-nil query result")
	}

	t.Logf("SUCCESS: Cross-network Iroh peer relay verified without direct IP!")
}

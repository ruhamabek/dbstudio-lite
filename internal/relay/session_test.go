package relay_test

import (
	"context"
	"testing"
	"time"

	"app-dev/internal/domain"
	"app-dev/internal/relay"
)

func TestHostSession_And_DialTicket(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	inspector := &FakeInspector{
		tables: []domain.Table{{Schema: "public", Name: "orders"}},
	}
	querier := &FakeQuerier{}

	host, ticketStr, err := relay.StartHostSession(ctx, "127.0.0.1:0", inspector, querier, relay.PermissionRead, "ecommerce_db")
	if err != nil {
		t.Fatalf("failed to start host session: %v", err)
	}
	defer host.Close()

	if ticketStr == "" {
		t.Fatal("expected non-empty ticket string")
	}

	client, err := relay.DialTicket(ctx, ticketStr)
	if err != nil {
		t.Fatalf("failed to dial ticket: %v", err)
	}
	defer client.Close()


	tables, err := client.ListTables(ctx)
	if err != nil {
		t.Fatalf("unexpected error listing tables: %v", err)
	}
	if len(tables) != 1 || tables[0].Name != "orders" {
		t.Fatalf("unexpected tables: %+v", tables)
	}


	result, err := client.Execute(ctx, "SELECT * FROM orders")
	if err != nil {
		t.Fatalf("unexpected query error: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil query result")
	}

	_, err = client.Execute(ctx, "DELETE FROM orders")
	if err == nil {
		t.Fatal("expected DELETE to be rejected under Read permission, got nil")
	}
}

func TestHostSession_RejectsInvalidToken(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	inspector := &FakeInspector{}
	querier := &FakeQuerier{}

	host, ticketStr, err := relay.StartHostSession(ctx, "127.0.0.1:0", inspector, querier, relay.PermissionRead, "testdb")
	if err != nil {
		t.Fatalf("failed to start host session: %v", err)
	}
	defer host.Close()

	parsed, err := relay.ParseTicket(ticketStr)
	if err != nil {
		t.Fatalf("failed to parse ticket: %v", err)
	}
	parsed.Token = "fake-hacker-token"
	tamperedTicket, _ := parsed.Encode()


	_, err = relay.DialTicket(ctx, tamperedTicket)
	if err == nil {
		t.Fatal("expected connection with tampered token to fail, got success")
	}
}
package relay_test

import (
	"testing"
	"app-dev/internal/relay"
)

func TestTicket_EncodeAndParse(t *testing.T) {

	t.Run("encodes and successfully parses a valid ticket", func(t *testing.T) {
	original := relay.Ticket{
		PeerID:  "12D3KooWDpJ7As7BWAwRMfu1VU2WCqNjvq387JEYKDBj4kx6nXTN",
		Addrs:  []string{"/ip4/127.0.0.1/tcp/4001", "/p2p-circuit"},
		Permission: relay.PermissionRead,
		Database: "testdb",
		Token:"secret-token-12345",
	}

	encoded, err := original.Encode()
	if err != nil {
		t.Fatalf("unexpected error duruing ecoding: %v", err)
	}

	if len(encoded) < 10 || encoded[:10] != "dbshare://" {
		t.Fatalf("excpected ticket to start with 'dbshare://', got %s", encoded)
	}

	parsed, err := relay.ParseTicket(encoded)
	if err != nil {
		t.Fatalf("unexpected error during parsing: %v", err)
	}

	if parsed.PeerID != original.PeerID {
		t.Errorf("expected PeerID %s, got %s", original.PeerID, parsed.PeerID)
	}

	
	if parsed.Permission != original.Permission {
		t.Errorf("got Permission %q, want %q", parsed.Permission, original.Permission)
	}

	if parsed.Database != original.Database {
		t.Errorf("got Database %q, want %q", parsed.Database, original.Database)
	}

	if parsed.Token != original.Token {
		t.Errorf("got Token %q, want %q", parsed.Token, original.Token)
	}

	if len(parsed.Addrs) != len(original.Addrs) {
		t.Errorf("got %d addrs, want %d", len(parsed.Addrs), len(original.Addrs))
	}

})

	t.Run("rejects ticket with invalid permision level", func(t *testing.T) {
        tkt := relay.Ticket{
			PeerID:     "some-peer",
			Permission: "admin_superuser_all_rights", // Not a valid Permission
			Token:      "token",
		}

		encoded, err := tkt.Encode()

		if err != nil {
			t.Fatalf("unexpected encode error: %v", err)
		}

		if _, err := relay.ParseTicket(encoded); err == nil {
			t.Fatal("expected error for unrecognized permission in ticket, got nil")
		}
	})
}
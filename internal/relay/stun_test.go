package relay_test

import (
	"context"
	"net"
	"testing"
	"time"

	"app-dev/internal/relay"
)

func TestSTUN_DiscoverPublicIP(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	ip, err := relay.DiscoverPublicIP(ctx, "stun.l.google.com:19302")
	if err != nil {
		t.Skipf("skipping live STUN test (no external network access): %v", err)
	}

	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		t.Fatalf("expected valid IP string from STUN, got: %q", ip)
	}

	if parsedIP.IsUnspecified() {
		t.Errorf("expected routable IP, got: %v", parsedIP)
	}
}
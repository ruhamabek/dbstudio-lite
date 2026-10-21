package relay_test

import (
	"context"
	"io"
	"net/netip"
	"testing"
	"time"

	"github.com/tmc/go-iroh/iroh"
	"github.com/tmc/go-iroh/netaddr"
)

func TestIrohStreamLoopback(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	const alpn = "dbstudio/test/1"

	server, err := iroh.Bind(ctx,
		iroh.WithALPNs(alpn),
		iroh.WithBindAddr(netip.AddrPortFrom(netip.IPv4Unspecified(), 0)),
	)
	if err != nil {
		t.Fatalf("failed to bind server: %v", err)
	}
	defer server.Shutdown(ctx)

	client, err := iroh.Bind(ctx,
		iroh.WithBindAddr(netip.AddrPortFrom(netip.IPv4Unspecified(), 0)),
	)
	if err != nil {
		t.Fatalf("failed to bind client: %v", err)
	}
	defer client.Shutdown(ctx)

	accepted := make(chan *iroh.Conn, 1)
	go func() {
		conn, err := server.Accept(ctx)
		if err != nil {
			t.Errorf("server.Accept failed: %v", err)
			return
		}
		accepted <- conn
	}()

	serverAddr := netaddr.NewEndpointAddr(server.ID()).WithIP(server.LocalAddr())
	clientConn, err := client.Connect(ctx, serverAddr, alpn)
	if err != nil {
		t.Fatalf("client.Connect failed: %v", err)
	}
	defer clientConn.CloseWithError(0, "")

	serverConn := <-accepted
	defer serverConn.CloseWithError(0, "")

	clientStream, err := clientConn.OpenStreamSync(ctx)
	if err != nil {
		t.Fatalf("client OpenStreamSync failed: %v", err)
	}
	defer clientStream.Close()

	msg := []byte("hello from dbstudio over iroh QUIC")

	
	go func() {
		_, _ = clientStream.Write(msg)
	}()

	serverStream, err := serverConn.AcceptStream(ctx)
	if err != nil {
		t.Fatalf("server AcceptStream failed: %v", err)
	}
	defer serverStream.Close()

	buf := make([]byte, len(msg))
	if _, err := io.ReadFull(serverStream, buf); err != nil {
		t.Fatalf("serverStream ReadFull failed: %v", err)
	}

	if string(buf) != string(msg) {
		t.Fatalf("received %q, want %q", string(buf), string(msg))
	}

	t.Logf("SUCCESS: Iroh peer connection and QUIC bidirectional stream verified!")
}

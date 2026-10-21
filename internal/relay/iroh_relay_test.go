package relay_test

import (
	"context"
	"io"
	"net"
	"testing"
	"time"

	"github.com/tmc/go-iroh/iroh"
	"github.com/tmc/go-iroh/relay"
)

func TestIrohPublicRelay(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping in short mode")
	}

	c, err := net.DialTimeout("tcp", "1.1.1.1:53", 500*time.Millisecond)
	if err != nil {
		t.Skip("skipping live public relay test: no internet connectivity")
	}
	c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	const alpn = "dbstudio/relay-test/1"

	server, err := iroh.Bind(ctx,
		iroh.WithALPNs(alpn),
		iroh.WithRelayMode(relay.ModeDefault()),
	)
	if err != nil {
		t.Fatalf("failed to bind server: %v", err)
	}
	defer server.Shutdown(ctx)

	
	if err := server.Online(ctx); err != nil {
		t.Fatalf("server.Online failed: %v", err)
	}

	addr := server.Addr()
	t.Logf("Server NodeID: %s, Relays: %v, Addrs: %v", server.ID(), addr.RelayURLs(), addr.IPAddrs())

	client, err := iroh.Bind(ctx,
		iroh.WithRelayMode(relay.ModeDefault()),
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

	clientConn, err := client.Connect(ctx, addr, alpn)
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

	msg := []byte("mentos-style iroh relay works!")
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
		t.Fatalf("got %q, want %q", string(buf), string(msg))
	}

	t.Logf("SUCCESS: Iroh cross-network connection over public n0 relay verified!")
}

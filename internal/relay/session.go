package relay

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/netip"
	"strings"
	"sync"
	"time"

	"app-dev/internal/domain"
	"github.com/tmc/go-iroh/iroh"
	"github.com/tmc/go-iroh/key"
	"github.com/tmc/go-iroh/netaddr"
	"github.com/tmc/go-iroh/relay"
)

const RelayALPN = "dbstudio/relay/1"

type HandshakeMsg struct {
	Token string `json:"token"`
}

type closerFunc func() error

func (f closerFunc) Close() error {
	if f == nil {
		return nil
	}
	return f()
}

type HostSession struct {
	listener     net.Listener
	irohEndpoint *iroh.Endpoint
	port         int
	token        string
	gateway      *HostGateway
	perm         Permission
	dbName       string
	closed       chan struct{}
	mu           sync.Mutex
	conns        []io.Closer
}

func getReachableAddrs(port int) []string {
	var addrs []string
	ifaces, err := net.Interfaces()
	if err == nil {
		for _, iface := range ifaces {
			if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
				continue
			}
			addrList, err := iface.Addrs()
			if err != nil {
				continue
			}
			for _, a := range addrList {
				if ipnet, ok := a.(*net.IPNet); ok && ipnet.IP.To4() != nil {
					addrs = append(addrs, fmt.Sprintf("%s:%d", ipnet.IP.String(), port))
				}
			}
		}
	}

	addrs = append(addrs, fmt.Sprintf("127.0.0.1:%d", port))
	return addrs
}

func StartHostSession(
	ctx context.Context,
	listenAddr string,
	inspector domain.SchemaInspector,
	querier domain.Querier,
	perm Permission,
	dbName string,
) (*HostSession, string, error) {
	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return nil, "", fmt.Errorf("failed to bind listener on %s: %w", listenAddr, err)
	}

	tokBytes := make([]byte, 16)
	if _, err := rand.Read(tokBytes); err != nil {
		ln.Close()
		return nil, "", err
	}
	token := hex.EncodeToString(tokBytes)

	tcpAddr, _ := ln.Addr().(*net.TCPAddr)
	port := tcpAddr.Port
	var addrs []string
	var irohEp *iroh.Endpoint
	var relayURLs []string
	var nodeID string

	if listenAddr == "127.0.0.1:0" || listenAddr == "localhost:0" {
		addrs = []string{ln.Addr().String()}
	} else {
		addrs = getReachableAddrs(port)

		// Mentos/Iroh cross-network peer reachability via n0 relay
		iCtx, iCancel := context.WithTimeout(ctx, 5*time.Second)
		ep, err := iroh.Bind(iCtx,
			iroh.WithALPNs(RelayALPN),
			iroh.WithRelayMode(relay.ModeDefault()),
		)
		iCancel()
		if err == nil {
			irohEp = ep
			nodeID = ep.ID().String()

			onlineCtx, oCancel := context.WithTimeout(ctx, 3*time.Second)
			_ = ep.Online(onlineCtx)
			oCancel()

			for _, ru := range ep.Addr().RelayURLs() {
				relayURLs = append(relayURLs, ru.String())
			}
			for _, ip := range ep.Addr().IPAddrs() {
				addrs = append(addrs, ip.String())
			}
		}
	}

	ticket := Ticket{
		PeerID:     "node-" + token[:8],
		Addrs:      addrs,
		NodeID:     nodeID,
		RelayURLs:  relayURLs,
		Permission: perm,
		Database:   dbName,
		Token:      token,
	}

	ticketStr, err := ticket.Encode()
	if err != nil {
		ln.Close()
		if irohEp != nil {
			_ = irohEp.Shutdown(context.Background())
		}
		return nil, "", fmt.Errorf("failed to encode share ticket: %w", err)
	}

	session := &HostSession{
		listener:     ln,
		irohEndpoint: irohEp,
		port:         port,
		token:        token,
		gateway:      NewHostGateway(inspector, querier),
		perm:         perm,
		dbName:       dbName,
		closed:       make(chan struct{}),
	}

	go session.acceptLoop(ctx)
	if irohEp != nil {
		go session.irohAcceptLoop(ctx)
	}

	return session, ticketStr, nil
}

func (s *HostSession) acceptLoop(ctx context.Context) {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			select {
			case <-s.closed:
				return
			case <-ctx.Done():
				return
			default:
				return
			}
		}

		s.mu.Lock()
		s.conns = append(s.conns, conn)
		s.mu.Unlock()

		go s.handleStream(ctx, conn)
	}
}

func (s *HostSession) irohAcceptLoop(ctx context.Context) {
	for {
		conn, err := s.irohEndpoint.Accept(ctx)
		if err != nil {
			select {
			case <-s.closed:
				return
			case <-ctx.Done():
				return
			default:
				return
			}
		}

		go s.handleIrohConn(ctx, conn)
	}
}

func (s *HostSession) handleIrohConn(ctx context.Context, conn *iroh.Conn) {
	defer conn.CloseWithError(0, "closed")
	for {
		stream, err := conn.AcceptStream(ctx)
		if err != nil {
			return
		}

		s.mu.Lock()
		s.conns = append(s.conns, stream)
		s.mu.Unlock()

		go s.handleStream(ctx, stream)
	}
}

func (s *HostSession) handleStream(ctx context.Context, rw io.ReadWriteCloser) {
	defer rw.Close()

	if c, ok := rw.(net.Conn); ok {
		_ = c.SetDeadline(time.Now().Add(5 * time.Second))
	}

	var hs HandshakeMsg
	if err := json.NewDecoder(rw).Decode(&hs); err != nil || hs.Token != s.token {
		_ = json.NewEncoder(rw).Encode(Response{Error: "unauthorized: invalid session token"})
		return
	}

	_ = json.NewEncoder(rw).Encode(Response{Result: json.RawMessage(`"ok"`)})

	if c, ok := rw.(net.Conn); ok {
		_ = c.SetDeadline(time.Time{})
	}

	_ = ServeStream(ctx, rw, s.gateway, s.perm)
}

func (s *HostSession) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	select {
	case <-s.closed:
		return nil
	default:
		close(s.closed)
	}

	if s.irohEndpoint != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		_ = s.irohEndpoint.Shutdown(ctx)
		cancel()
	}

	for _, c := range s.conns {
		_ = c.Close()
	}
	return s.listener.Close()
}

func DialTicket(ctx context.Context, ticketStr string) (*RelayClient, error) {
	ticket, err := ParseTicket(ticketStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ticket: %w", err)
	}

	var stream io.ReadWriteCloser
	var closeFunc func() error
	var dialErrors []string

	if len(ticket.Addrs) > 0 {
		dialer := net.Dialer{Timeout: 1500 * time.Millisecond}
		for _, addr := range ticket.Addrs {
			c, err := dialer.DialContext(ctx, "tcp", addr)
			if err == nil {
				stream = c
				closeFunc = c.Close
				break
			}
			dialErrors = append(dialErrors, fmt.Sprintf("tcp://%s: %v", addr, err))
		}
	}

	if stream == nil && ticket.NodeID != "" {
		ep, err := iroh.Bind(ctx, iroh.WithRelayMode(relay.ModeDefault()))
		if err != nil {
			dialErrors = append(dialErrors, fmt.Sprintf("iroh bind failed: %v", err))
		} else {
			nodeID, err := key.ParseEndpointID(ticket.NodeID)
			if err != nil {
				_ = ep.Shutdown(ctx)
				dialErrors = append(dialErrors, fmt.Sprintf("invalid ticket nodeID %q: %v", ticket.NodeID, err))
			} else {
				addr := netaddr.NewEndpointAddr(nodeID)
				for _, r := range ticket.RelayURLs {
					if ru, err := netaddr.ParseRelayURL(r); err == nil {
						addr = addr.WithRelayURL(ru)
					}
				}
				for _, a := range ticket.Addrs {
					if ap, err := netip.ParseAddrPort(a); err == nil {
						addr = addr.WithIP(ap)
					}
				}

				connCtx, cCancel := context.WithTimeout(ctx, 15*time.Second)
				conn, err := ep.Connect(connCtx, addr, RelayALPN)
				cCancel()
				if err != nil {
					_ = ep.Shutdown(ctx)
					dialErrors = append(dialErrors, fmt.Sprintf("iroh connect to %s failed: %v", ticket.NodeID, err))
				} else {
					s, err := conn.OpenStreamSync(ctx)
					if err != nil {
						conn.CloseWithError(1, "stream error")
						_ = ep.Shutdown(ctx)
						dialErrors = append(dialErrors, fmt.Sprintf("iroh open stream failed: %v", err))
					} else {
						stream = s
						closeFunc = func() error {
							_ = s.Close()
							conn.CloseWithError(0, "done")
							return ep.Shutdown(context.Background())
						}
					}
				}
			}
		}
	}

	if stream == nil {
		if len(dialErrors) > 0 {
			return nil, fmt.Errorf("could not connect to peer: %s", strings.Join(dialErrors, "; "))
		}
		return nil, errors.New("ticket contains no reachable network addresses or relay identifiers")
	}

	if c, ok := stream.(net.Conn); ok {
		_ = c.SetDeadline(time.Now().Add(5 * time.Second))
	}

	if err := json.NewEncoder(stream).Encode(HandshakeMsg{Token: ticket.Token}); err != nil {
		if closeFunc != nil {
			_ = closeFunc()
		}
		return nil, fmt.Errorf("handshake write failed: %w", err)
	}

	var hsResp Response
	if err := json.NewDecoder(stream).Decode(&hsResp); err != nil {
		if closeFunc != nil {
			_ = closeFunc()
		}
		return nil, fmt.Errorf("handshake read failed: %w", err)
	}
	if hsResp.Error != "" {
		if closeFunc != nil {
			_ = closeFunc()
		}
		return nil, errors.New(hsResp.Error)
	}

	if c, ok := stream.(net.Conn); ok {
		_ = c.SetDeadline(time.Time{})
	}

	caller := NewStreamCaller(stream)
	client := NewRelayClient(caller)
	if closeFunc != nil {
		client.closer = closerFunc(closeFunc)
	}
	return client, nil
}

package relay

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"app-dev/internal/domain"
)

type HandshakeMsg struct {
	Token string `json:"token"`
}

type HostSession struct {
	listener net.Listener
	port     int
	upnp     *UPnPClient
	token    string
	gateway  *HostGateway
	perm     Permission
	dbName   string
	closed   chan struct{}
	mu       sync.Mutex
	conns    []net.Conn
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
	var upnpClient *UPnPClient

	if listenAddr == "127.0.0.1:0" || listenAddr == "localhost:0" {
		addrs = []string{ln.Addr().String()}
	} else {
		addrs = getReachableAddrs(port)

 		stunCtx, cancel := context.WithTimeout(ctx, 1500*time.Millisecond)
		if pubIP, err := DiscoverPublicIP(stunCtx, ""); err == nil && pubIP != "" {
			addrs = append([]string{fmt.Sprintf("%s:%d", pubIP, port)}, addrs...)
		}
		cancel()

 		upnpClient = NewUPnPClient()
		localIP := "127.0.0.1"
		for _, a := range addrs {
			if h, _, err := net.SplitHostPort(a); err == nil && !strings.HasPrefix(h, "127.") {
				localIP = h
				break
			}
		}
		upnpCtx, uCancel := context.WithTimeout(ctx, 1*time.Second)
		_ = upnpClient.ForwardPort(upnpCtx, port, localIP)
		uCancel()
	}

	ticket := Ticket{
		PeerID:     "node-" + token[:8],
		Addrs:      addrs,
		Permission: perm,
		Database:   dbName,
		Token:      token,
	}

	ticketStr, err := ticket.Encode()
	if err != nil {
		ln.Close()
		return nil, "", fmt.Errorf("failed to encode share ticket: %w", err)
	}

	session := &HostSession{
		listener: ln,
		port:     port,
		upnp:     upnpClient,
		token:    token,
		gateway:  NewHostGateway(inspector, querier),
		perm:     perm,
		dbName:   dbName,
		closed:   make(chan struct{}),
	}

	go session.acceptLoop(ctx)

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

		go s.handleConn(ctx, conn)
	}
}

func (s *HostSession) handleConn(ctx context.Context, conn net.Conn) {
	defer conn.Close()

	_ = conn.SetDeadline(time.Now().Add(3 * time.Second))
	var hs HandshakeMsg
	if err := json.NewDecoder(conn).Decode(&hs); err != nil || hs.Token != s.token {
		_ = json.NewEncoder(conn).Encode(Response{Error: "unauthorized: invalid session token"})
		return
	}

	_ = json.NewEncoder(conn).Encode(Response{Result: json.RawMessage(`"ok"`)})
	_ = conn.SetDeadline(time.Time{})

	_ = ServeStream(ctx, conn, s.gateway, s.perm)
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

	// Clean up router UPnP port mapping
	if s.upnp != nil && s.port > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		_ = s.upnp.DeletePort(ctx, s.port)
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

	if len(ticket.Addrs) == 0 {
		return nil, errors.New("ticket contains no network addresses")
	}

	var conn net.Conn
	var dialErr error

	dialer := net.Dialer{Timeout: 3 * time.Second}
	for _, addr := range ticket.Addrs {
		c, err := dialer.DialContext(ctx, "tcp", addr)
		if err == nil {
			conn = c
			break
		}
		dialErr = err
	}

	if conn == nil {
		return nil, fmt.Errorf("could not connect to peer: %w", dialErr)
	}

	_ = conn.SetDeadline(time.Now().Add(3 * time.Second))
	if err := json.NewEncoder(conn).Encode(HandshakeMsg{Token: ticket.Token}); err != nil {
		conn.Close()
		return nil, fmt.Errorf("handshake write failed: %w", err)
	}

	var hsResp Response
	if err := json.NewDecoder(conn).Decode(&hsResp); err != nil {
		conn.Close()
		return nil, fmt.Errorf("handshake read failed: %w", err)
	}
	if hsResp.Error != "" {
		conn.Close()
		return nil, errors.New(hsResp.Error)
	}
	_ = conn.SetDeadline(time.Time{})

	caller := NewStreamCaller(conn)
	return NewRelayClient(caller), nil
}


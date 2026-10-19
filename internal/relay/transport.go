package relay

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
)

type StreamCaller struct {
	mu  sync.Mutex
	rw  io.ReadWriter 
	enc *json.Encoder
	dec *json.Decoder
}
func NewStreamCaller(rw io.ReadWriter) *StreamCaller {
	return &StreamCaller{
		rw:  rw, 
		enc: json.NewEncoder(rw),
		dec: json.NewDecoder(rw),
	}
}
func (c *StreamCaller) Call(ctx context.Context, req Request) (Response, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := c.enc.Encode(req); err != nil {
		return Response{}, fmt.Errorf("failed to encode request across stream: %w", err)
	}

	var resp Response
	if err := c.dec.Decode(&resp); err != nil {
		return Response{}, fmt.Errorf("failed to decode response from stream: %w", err)
	}

	return resp, nil
}

func (c *StreamCaller) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if closer, ok := c.rw.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

func ServeStream(ctx context.Context, rw io.ReadWriter, gateway *HostGateway, perm Permission) error {
	dec := json.NewDecoder(rw)
	enc := json.NewEncoder(rw)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		var req Request
		if err := dec.Decode(&req); err != nil {
			if errors.Is(err, io.EOF) || errors.Is(err, net.ErrClosed) {
				return nil
			}
			return err
		}

		resp := gateway.Handle(ctx, perm, req)
		if err := enc.Encode(resp); err != nil {
			return err
		}
	}
}
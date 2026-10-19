package relay

import (
	"app-dev/internal/domain"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

type Caller interface {
	Call(ctx context.Context, req Request) (Response, error)
} 

type RelayClient struct {
	caller Caller
	closer io.Closer
}

func NewRelayClient(caller Caller) *RelayClient {
	return &RelayClient{caller: caller}
}

func (c *RelayClient) Close() error {
	if c.closer != nil {
		return c.closer.Close()
	}
	return nil
}

func (c *RelayClient) nextID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func decodeResult[T any](raw json.RawMessage) (T, error) {
    var t T
    if len(raw) == 0 {
        return t, nil
    }

    if err := json.Unmarshal(raw, &t); err != nil {
        return t, fmt.Errorf("failed to unmarshal into target type: %w", err)
    }

    return t, nil
}

func (c *RelayClient) ListTables(ctx context.Context) ([]domain.Table, error) {
	resp, err := c.caller.Call(ctx, Request{
		ID: c.nextID(),
		Method: MethodListTables,
	})

	if err != nil {
		 return nil, err
	}
	
	if resp.Error != "" {
		return nil, errors.New(resp.Error)
	}
	
	return decodeResult[[]domain.Table](resp.Result)
}

func (c *RelayClient) GetColumns(ctx context.Context, schema, table string) ([]domain.Column, error) {
	resp, err := c.caller.Call(ctx, Request{
		ID:     c.nextID(),
		Method: MethodGetColumns,
		Schema: schema,
		Table:  table,
	})
	if err != nil {
		return nil, err
	}
	if resp.Error != "" {
		return nil, errors.New(resp.Error)
	}
	return decodeResult[[]domain.Column](resp.Result)
}

func (c *RelayClient) GetPrimaryKeys(ctx context.Context, schema, table string) ([]string, error) {
	resp, err := c.caller.Call(ctx, Request{
		ID:     c.nextID(),
		Method: MethodGetPrimaryKeys,
		Schema: schema,
		Table:  table,
	})
	if err != nil {
		return nil, err
	}
	if resp.Error != "" {
		return nil, errors.New(resp.Error)
	}
	return decodeResult[[]string](resp.Result)
}


func (c *RelayClient) Execute(ctx context.Context, sql string) (*domain.QueryResult, error) {
	resp, err := c.caller.Call(ctx, Request{
		ID:     c.nextID(),
		Method: MethodExecuteQuery,
		SQL:    sql,
	})
	if err != nil {
		return nil, err
	}
	if resp.Error != "" {
		return nil, errors.New(resp.Error)
	}
	return decodeResult[*domain.QueryResult](resp.Result)
}
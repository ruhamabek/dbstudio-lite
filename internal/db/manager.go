package db

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

var ErrNotConnected = errors.New("no active database connection")

type Client interface {
	Querier
	Ping(ctx context.Context) error
	Close() error
}

type Connector func(ctx context.Context, dsn string) (Client, error)

type ConnectionManager struct {
	mu        sync.RWMutex
	connector Connector
	client    Client
	inspector *Inspector
	runner    *Runner
}

func NewConnectionManager(connector Connector) *ConnectionManager {
	return &ConnectionManager{
		connector: connector,
	}
}


func (m *ConnectionManager) Inspector() (*Inspector, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.client == nil {
		return nil, ErrNotConnected
	}
	return m.inspector, nil
}

func (m *ConnectionManager) Runner() (*Runner, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.client == nil {
		return nil, ErrNotConnected
	}
	return m.runner, nil
}

func (m *ConnectionManager) Connect(ctx context.Context, cfg Config) error{
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.connector == nil {
		return errors.New("No connector configured")
	}

	if m.client != nil {
		_ = m.client.Close()
		m.client = nil
		m.inspector = nil
		m.runner = nil
	}

	dsn := cfg.ConnectionString()

	client, err := m.connector(ctx, dsn)

	if err != nil {
		return fmt.Errorf("Connection Failed: %v", err)
	}

	m.client = client
	m.inspector = NewInspector(client)
	m.runner = NewRunner(client)

	return nil
}

func (m *ConnectionManager) Disconnect() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.client == nil {
		return nil
	}

	err := m.client.Close()
	m.client = nil 
	m.inspector = nil 
	m.runner =nil 

	if err != nil {
		return fmt.Errorf("failed to close connection: %v", err)
	}

	return nil
}

func (m *ConnectionManager) Ping(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.client == nil {
		return ErrNotConnected
	}

	return m.client.Ping(ctx)
}
package service

import (
	"app-dev/internal/domain"
	"app-dev/internal/driver"
	"context"
	"errors"
	"fmt"
	"sync"
)

var ErrNotConnected = errors.New("no active database connection")

type ConnectionManager struct {
	mu        sync.RWMutex
	connector domain.Connector
	client    domain.Client
	inspector domain.SchemaInspector
	runner    *Runner
}

func NewConnectionManager(connector domain.Connector) *ConnectionManager {
	return &ConnectionManager{
		connector: connector,
	}
}

func (m *ConnectionManager) Inspector() (domain.SchemaInspector, error) {
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

func (m *ConnectionManager) Connect(ctx context.Context, cfg domain.Config) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.connector == nil {
		return errors.New("no connector configured")
	}

	if m.client != nil {
		_ = m.client.Close()
		m.client = nil
		m.inspector = nil
		m.runner = nil
	}

	client, err := m.connector(ctx, cfg)
	if err != nil {
		return fmt.Errorf("connection failed: %w", err)
	}

	m.client = client
	m.inspector = driver.NewInspector(cfg.Driver, client)
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
	m.runner = nil

	if err != nil {
		return fmt.Errorf("failed to close connection: %w", err)
	}

	return nil
}

func (m *ConnectionManager) Ping(ctx context.Context) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.client == nil {
		return ErrNotConnected
	}

	return m.client.Ping(ctx)
}

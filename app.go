package main

import (
	"app-dev/internal/domain"
	"app-dev/internal/relay"
	"app-dev/internal/service"
	"context"
	"fmt"
)

type App struct {
	ctx          context.Context
	manager      *service.ConnectionManager
	store        domain.ConnectionStore
	relaySession *relay.HostSession
	relayClient  *relay.RelayClient
}

func NewApp(manager *service.ConnectionManager, store domain.ConnectionStore) *App {
	return &App{
		manager: manager,
		store:   store,
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) Connect(cfg domain.Config) error {
	return a.manager.Connect(a.ctx, cfg)
}

func (a *App) ConnectWithURL(rawURL string) error {
	cfg, err := domain.ParseURL(rawURL)
	if err != nil {
		return err
	}
	return a.manager.Connect(a.ctx, cfg)
}

func (a *App) Disconnect() error {
	if a.relayClient != nil {
		_ = a.relayClient.Close()
		a.relayClient = nil
	}
	return a.manager.Disconnect()
}

func (a *App) Ping() error {
	if a.relayClient != nil {
		return nil
	}
	return a.manager.Ping(a.ctx)
}

func (a *App) GetTables() ([]domain.Table, error) {
	if a.relayClient != nil {
		return a.relayClient.ListTables(a.ctx)
	}
	inspector, err := a.manager.Inspector()
	if err != nil {
		return nil, err
	}
	return inspector.ListTables(a.ctx)
}

func (a *App) GetColumns(schema, table string) ([]domain.Column, error) {
	if a.relayClient != nil {
		return a.relayClient.GetColumns(a.ctx, schema, table)
	}
	inspector, err := a.manager.Inspector()
	if err != nil {
		return nil, err
	}
	return inspector.GetColumns(a.ctx, schema, table)
}

func (a *App) GetPrimaryKeys(schema, table string) ([]string, error) {
	if a.relayClient != nil {
		return a.relayClient.GetPrimaryKeys(a.ctx, schema, table)
	}
	inspector, err := a.manager.Inspector()
	if err != nil {
		return nil, err
	}
	return inspector.GetPrimaryKeys(a.ctx, schema, table)
}

func (a *App) ExecuteQuery(sql string) (*domain.QueryResult, error) {
	if a.relayClient != nil {
		return a.relayClient.Execute(a.ctx, sql)
	}
	runner, err := a.manager.Runner()
	if err != nil {
		return nil, err
	}
	return runner.Execute(a.ctx, sql)
}

func (a *App) GetSavedConnections() ([]domain.SavedConnection, error) {
	if a.store == nil {
		return []domain.SavedConnection{}, nil
	}
	return a.store.List()
}

func (a *App) SaveConnection(conn domain.SavedConnection) (domain.SavedConnection, error) {
	if a.store == nil {
		return domain.SavedConnection{}, fmt.Errorf("no store configured")
	}
	return a.store.Save(conn)
}

func (a *App) DeleteConnection(id string) error {
	if a.store == nil {
		return nil
	}
	return a.store.Delete(id)
}

func (a *App) StartSharing(perm string) (string, error) {
	if a.relaySession != nil {
		_ = a.relaySession.Close()
		a.relaySession = nil
	}

	inspector, err := a.manager.Inspector()
	if err != nil {
		return "", fmt.Errorf("must connect to a database before sharing: %w", err)
	}

	querier, err := a.manager.Querier()
	if err != nil {
		return "", fmt.Errorf("must connect to a database before sharing: %w", err)
	}

	session, ticketStr, err := relay.StartHostSession(
		a.ctx,
		":0",
		inspector,
		querier,
		relay.Permission(perm),
		"shared-db",
	)
	if err != nil {
		return "", err
	}

	a.relaySession = session
	return ticketStr, nil
}

func (a *App) StopSharing() error {
	if a.relaySession != nil {
		err := a.relaySession.Close()
		a.relaySession = nil
		return err
	}
	return nil
}

func (a *App) ConnectRelay(ticketStr string) error {
	client, err := relay.DialTicket(a.ctx, ticketStr)
	if err != nil {
		return err
	}

	_ = a.Disconnect()
	a.relayClient = client
	return nil
}
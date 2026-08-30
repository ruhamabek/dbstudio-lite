package main

import (
	"app-dev/internal/db"
	"context"
	"fmt"
)


type App struct {
	ctx     context.Context
	manager *db.ConnectionManager
	store   db.ConnectionStore
}

func NewApp(manager *db.ConnectionManager, store db.ConnectionStore) *App {
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

func (a *App) Connect(cfg db.Config) error {
	return a.manager.Connect(a.ctx, cfg)
}

func (a *App) ConnectWithURL(rawURL string) error {
	cfg, err := db.ParseURL(rawURL)
	if err != nil {
		return err
	}
	return a.manager.Connect(a.ctx, cfg)
}

func (a *App) Disconnect() error {
	return a.manager.Disconnect()
}

func (a *App) Ping() error {
	return a.manager.Ping(a.ctx)
}

func (a *App) GetTables() ([]db.Table, error) {
	inspector, err := a.manager.Inspector()
	if err != nil {
		return nil, err
	}
	return inspector.ListTables(a.ctx)
}

func (a *App) GetColumns(schema, table string) ([]db.Column, error) {
	inspector, err := a.manager.Inspector()
	if err != nil {
		return nil, err
	}
	return inspector.GetColumns(a.ctx, schema, table)
}


func (a *App) GetPrimaryKeys(schema, table string) ([]string, error) {
	inspector, err := a.manager.Inspector()
	if err != nil {
		return nil, err
	}
	return inspector.GetPrimaryKeys(a.ctx, schema, table)
}

func (a *App) ExecuteQuery(sql string) (*db.QueryResult, error) {
	runner, err := a.manager.Runner()
	if err != nil {
		return nil, err
	}
	return runner.Execute(a.ctx, sql)
}

func (a *App) GetSavedConnections() ([]db.SavedConnection, error) {
	if a.store == nil {
		return []db.SavedConnection{}, nil
	}
	return a.store.List()
}

func (a *App) SaveConnection(conn db.SavedConnection) (db.SavedConnection, error) {
	if a.store == nil {
		return db.SavedConnection{}, fmt.Errorf("no store configured")
	}
	return a.store.Save(conn)
}

func (a *App) DeleteConnection(id string) error {
	if a.store == nil {
		return nil
	}
	return a.store.Delete(id)
}
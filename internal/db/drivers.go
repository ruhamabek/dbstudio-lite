package db

import (
	"context"
 )

 

 func NewMultiDriverConnector() Connector {
	return func(ctx context.Context, cfg Config) (Client, error) {
		switch cfg.Driver {
		case DriverMySQL:
			return NewMySQLClient(ctx, cfg.ConnectionString())
		case DriverSQLite:
			return NewSQLiteClient(ctx, cfg.ConnectionString())
		case DriverPostgres:
			return NewPgxClient(ctx, cfg.ConnectionString())
		default:
 			return NewPgxClient(ctx, cfg.ConnectionString())
		}
	}
}
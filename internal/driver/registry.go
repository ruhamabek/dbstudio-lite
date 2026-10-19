package driver

import (
	"app-dev/internal/domain"
	"app-dev/internal/driver/mysql"
	"app-dev/internal/driver/postgres"
	"app-dev/internal/driver/sqlite"
	"context"
)

// NewMultiDriverConnector returns a Connector that dispatches to Postgres, MySQL, or SQLite.
func NewMultiDriverConnector() domain.Connector {
	return func(ctx context.Context, cfg domain.Config) (domain.Client, error) {
		switch cfg.Driver {
		case domain.DriverMySQL:
			return mysql.NewClient(ctx, cfg.ConnectionString())
		case domain.DriverSQLite:
			return sqlite.NewClient(ctx, cfg.ConnectionString())
		case domain.DriverPostgres:
			return postgres.NewClient(ctx, cfg.ConnectionString())
		default:
			return postgres.NewClient(ctx, cfg.ConnectionString())
		}
	}
}

// NewInspector constructs the appropriate domain.SchemaInspector based on DriverType.
func NewInspector(driver domain.DriverType, db domain.Querier) domain.SchemaInspector {
	switch driver {
	case domain.DriverMySQL:
		return mysql.NewInspector(db)
	case domain.DriverSQLite:
		return sqlite.NewInspector(db)
	default:
		return postgres.NewInspector(db)
	}
}

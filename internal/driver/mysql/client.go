package mysql

import (
	"app-dev/internal/domain"
	"app-dev/internal/driver/sqlcommon"
	"context"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func NewClient(ctx context.Context, dsn string) (domain.Client, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open mysql connection: %w", err)
	}

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("failed to ping mysql database: %w", err)
	}

	return &sqlcommon.SQLClient{DB: db}, nil
}
